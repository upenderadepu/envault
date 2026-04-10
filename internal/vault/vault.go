package vault

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bhartiyaanshul/envault/internal/config"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var vaultOperationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "envault_vault_operations_total",
	Help: "Total number of Vault operations.",
}, []string{"operation", "status"})

type VaultService struct {
	client      *vault.Client
	mountPrefix string
}

func NewVaultService(cfg config.VaultConfig) (*VaultService, error) {
	client, err := vault.New(
		vault.WithAddress(cfg.Addr),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("vault client creation failed: %w", err)
	}

	if err := client.SetToken(cfg.Token); err != nil {
		return nil, fmt.Errorf("vault set token failed: %w", err)
	}

	log.Info().Str("addr", cfg.Addr).Msg("vault client connected")
	return &VaultService{client: client, mountPrefix: cfg.MountPrefix}, nil
}

// StartTokenRenewal runs a background goroutine that renews the server's
// Vault token every 30 minutes. Stops when ctx is canceled.
func (v *VaultService) StartTokenRenewal(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("vault token renewal stopped")
				return
			case <-ticker.C:
				_, err := v.client.Auth.TokenRenewSelf(ctx, schema.TokenRenewSelfRequest{})
				if err != nil {
					log.Error().Err(err).Msg("vault token renewal failed")
				} else {
					log.Debug().Msg("vault token renewed")
				}
			}
		}
	}()
}

// EnableKV2 enables a KV-v2 secrets engine at the given mount path.
// Mount creation is NOT idempotent — checks if it exists first.
func (v *VaultService) EnableKV2(ctx context.Context, mountPath string) error {
	// Check if mount already exists
	mounts, err := v.client.System.MountsListSecretsEngines(ctx)
	if err != nil {
		return fmt.Errorf("failed to list mounts: %w", err)
	}

	mountKey := mountPath + "/"
	if mounts.Data != nil {
		if _, exists := mounts.Data[mountKey]; exists {
			log.Debug().Str("mount", mountPath).Msg("KV-v2 mount already exists")
			return nil
		}
	}

	_, err = v.client.System.MountsEnableSecretsEngine(ctx, mountPath, schema.MountsEnableSecretsEngineRequest{
		Type: "kv",
		Options: map[string]interface{}{
			"version": "2",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to enable KV-v2 at %s: %w", mountPath, err)
	}

	log.Info().Str("mount", mountPath).Msg("KV-v2 secrets engine enabled")
	return nil
}

// WritePolicy writes an HCL policy to Vault.
func (v *VaultService) WritePolicy(ctx context.Context, policyName, policyHCL string) error {
	_, err := v.client.System.PoliciesWriteAclPolicy(ctx, policyName, schema.PoliciesWriteAclPolicyRequest{
		Policy: policyHCL,
	})
	if err != nil {
		return fmt.Errorf("failed to write policy %s: %w", policyName, err)
	}
	log.Debug().Str("policy", policyName).Msg("vault policy written")
	return nil
}

// BuildPolicies generates HCL policy string for a given role.
//
// Roles:
//   - admin: full access to all environment paths
//   - developer: read/write to development and staging, no production
//   - ci: read-only to production and staging
func (v *VaultService) BuildPolicies(mountPath, role string, environments []string) string {
	var b strings.Builder

	switch role {
	case "admin":
		fmt.Fprintf(&b, "path \"%s/data/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n\n", mountPath)
		fmt.Fprintf(&b, "path \"%s/metadata/*\" {\n  capabilities = [\"read\", \"list\", \"delete\"]\n}\n", mountPath)

	case "developer":
		for _, env := range environments {
			if env == "production" {
				continue
			}
			fmt.Fprintf(&b, "path \"%s/data/%s\" {\n  capabilities = [\"create\", \"read\", \"update\", \"list\"]\n}\n\n", mountPath, env)
			fmt.Fprintf(&b, "path \"%s/metadata/%s\" {\n  capabilities = [\"read\", \"list\"]\n}\n\n", mountPath, env)
		}

	case "ci":
		for _, env := range environments {
			if env == "development" {
				continue
			}
			fmt.Fprintf(&b, "path \"%s/data/%s\" {\n  capabilities = [\"read\"]\n}\n\n", mountPath, env)
			fmt.Fprintf(&b, "path \"%s/metadata/%s\" {\n  capabilities = [\"read\", \"list\"]\n}\n\n", mountPath, env)
		}
	}

	return b.String()
}

// CreateUserToken creates a new Vault token with the given policies and TTL.
// Returns the token (given to user ONCE) and the accessor (stored in DB).
func (v *VaultService) CreateUserToken(ctx context.Context, policies []string, ttl, maxTTL time.Duration) (token string, accessor string, err error) {
	resp, err := v.client.Auth.TokenCreate(ctx, schema.TokenCreateRequest{
		Policies:       policies,
		Ttl:            ttl.String(),
		ExplicitMaxTtl: maxTTL.String(),
		Renewable:      true,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create token: %w", err)
	}

	vaultOperationsTotal.WithLabelValues("create_token", "success").Inc()
	return resp.Auth.ClientToken, resp.Auth.Accessor, nil
}

// RevokeTokenByAccessor revokes a Vault token using only its accessor.
// This is how we revoke tokens without ever storing the token itself.
func (v *VaultService) RevokeTokenByAccessor(ctx context.Context, accessor string) error {
	if accessor == "" {
		return nil
	}
	_, err := v.client.Auth.TokenRevokeAccessor(ctx, schema.TokenRevokeAccessorRequest{
		Accessor: accessor,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke token by accessor: %w", err)
	}
	return nil
}

// RevokeAllProjectCredentials revokes all tokens for a project (on project delete).
func (v *VaultService) RevokeAllProjectCredentials(ctx context.Context, accessors []string) error {
	var errs []string
	for _, acc := range accessors {
		if err := v.RevokeTokenByAccessor(ctx, acc); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to revoke %d tokens: %s", len(errs), strings.Join(errs, "; "))
	}
	return nil
}

// ReadSecret reads all key-value pairs at a KV-v2 path.
// Returns empty map if path has no data yet (not an error).
func (v *VaultService) ReadSecret(ctx context.Context, mountPath, secretPath string) (map[string]interface{}, error) {
	resp, err := v.client.Secrets.KvV2Read(ctx, secretPath, vault.WithMountPath(mountPath))
	if err != nil {
		// Vault returns 404 if path doesn't exist yet — treat as empty
		if strings.Contains(err.Error(), "404") {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("failed to read secret at %s/%s: %w", mountPath, secretPath, err)
	}

	if resp.Data.Data == nil {
		return make(map[string]interface{}), nil
	}
	vaultOperationsTotal.WithLabelValues("read", "success").Inc()
	return resp.Data.Data, nil
}

// WriteSecret performs a read-merge-write to avoid overwriting existing keys.
// New keys overwrite existing ones; keys not in `data` are preserved.
func (v *VaultService) WriteSecret(ctx context.Context, mountPath, secretPath string, data map[string]interface{}) error {
	existing, err := v.ReadSecret(ctx, mountPath, secretPath)
	if err != nil {
		return err
	}

	// Merge: new data overwrites existing keys, preserves others
	for k, val := range data {
		existing[k] = val
	}

	_, err = v.client.Secrets.KvV2Write(ctx, secretPath, schema.KvV2WriteRequest{
		Data: existing,
	}, vault.WithMountPath(mountPath))
	if err != nil {
		vaultOperationsTotal.WithLabelValues("write", "error").Inc()
		return fmt.Errorf("failed to write secret at %s/%s: %w", mountPath, secretPath, err)
	}
	vaultOperationsTotal.WithLabelValues("write", "success").Inc()
	return nil
}

// DeleteSecretKey removes a single key from a KV-v2 path.
// Reads current data, removes the key, writes back.
func (v *VaultService) DeleteSecretKey(ctx context.Context, mountPath, secretPath, key string) error {
	existing, err := v.ReadSecret(ctx, mountPath, secretPath)
	if err != nil {
		return err
	}

	delete(existing, key)

	_, err = v.client.Secrets.KvV2Write(ctx, secretPath, schema.KvV2WriteRequest{
		Data: existing,
	}, vault.WithMountPath(mountPath))
	if err != nil {
		vaultOperationsTotal.WithLabelValues("delete", "error").Inc()
		return fmt.Errorf("failed to write after deleting key %s: %w", key, err)
	}
	vaultOperationsTotal.WithLabelValues("delete", "success").Inc()
	return nil
}
