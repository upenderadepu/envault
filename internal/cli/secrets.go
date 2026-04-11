package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"s"},
	Short:   "Manage secrets",
	Run:     secretsListRun,
}

func init() {
	secretsCmd.PersistentFlags().String("env", "development", "Environment (development, staging, production)")
	secretsCmd.AddCommand(secretsSetCmd)
	secretsCmd.AddCommand(secretsGetCmd)
	secretsCmd.AddCommand(secretsDeleteCmd)
}

func secretsListRun(cmd *cobra.Command, args []string) {
	slug := getProjectSlug()
	env := getEnv(cmd)
	client := newAuthClient()

	sp := startSpinner(fmt.Sprintf("Loading %s secrets...", env))
	body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/secrets?environment=%s", slug, env), nil)
	if err != nil {
		sp.fail("Failed: " + err.Error())
		os.Exit(1)
	}
	sp.stop(fmt.Sprintf("Secrets loaded (%s)", env))
	fmt.Println()

	var secrets []struct {
		KeyName      string `json:"key_name"`
		VaultVersion int    `json:"vault_version"`
		LastModified string `json:"last_modified_at"`
	}
	json.Unmarshal(body, &secrets)

	if len(secrets) == 0 {
		info("No secrets in '%s'. Add one with 'envault secrets set KEY=VALUE'", env)
		fmt.Println()
		return
	}

	rows := make([][]string, len(secrets))
	for i, s := range secrets {
		modified := s.LastModified
		if len(modified) > 10 {
			modified = modified[:10]
		}
		rows[i] = []string{s.KeyName, fmt.Sprintf("v%d", s.VaultVersion), modified}
	}
	printTable([]string{"Key", "Version", "Modified"}, rows)
	fmt.Println()
	dim.Printf("  %d secret(s) in %s. Use 'envault secrets get <KEY>' to reveal a value.\n\n", len(secrets), env)
}

var secretsSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE [KEY2=VALUE2 ...]",
	Short: "Set one or more secrets",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		env := getEnv(cmd)
		client := newAuthClient()

		if len(args) == 1 {
			// Single secret
			key, value, ok := parseKeyValue(args[0])
			if !ok {
				fatal("Invalid format. Use KEY=VALUE")
			}

			sp := startSpinner(fmt.Sprintf("Setting %s...", key))
			_, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/secrets", slug), map[string]string{
				"environment": env,
				"key":         key,
				"value":       value,
			})
			if err != nil {
				sp.fail("Failed: " + err.Error())
				os.Exit(1)
			}
			sp.stop(fmt.Sprintf("Set %s in %s", key, env))
		} else {
			// Bulk set
			secrets := make(map[string]string)
			for _, arg := range args {
				key, value, ok := parseKeyValue(arg)
				if !ok {
					fatal("Invalid format '%s'. Use KEY=VALUE", arg)
				}
				secrets[key] = value
			}

			sp := startSpinner(fmt.Sprintf("Setting %d secrets...", len(secrets)))
			_, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/secrets/bulk", slug), map[string]interface{}{
				"environment": env,
				"secrets":     secrets,
			})
			if err != nil {
				sp.fail("Failed: " + err.Error())
				os.Exit(1)
			}
			sp.stop(fmt.Sprintf("Set %d secrets in %s", len(secrets), env))
		}
	},
}

var secretsGetCmd = &cobra.Command{
	Use:   "get <KEY>",
	Short: "Get a secret value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		env := getEnv(cmd)
		key := args[0]
		client := newAuthClient()

		valueOnly, _ := cmd.Flags().GetBool("value-only")

		body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/secrets/%s?environment=%s", slug, key, env), nil)
		if err != nil {
			fatal("%s", err.Error())
		}

		var resp struct {
			Value string `json:"value"`
		}
		json.Unmarshal(body, &resp)

		if valueOnly {
			fmt.Print(resp.Value)
		} else {
			fmt.Println()
			printKeyValue("Key", key)
			printKeyValue("Env", env)
			printKeyValue("Value", resp.Value)
			fmt.Println()
		}
	},
}

func init() {
	secretsGetCmd.Flags().Bool("value-only", false, "Print only the value (for scripting)")
}

var secretsDeleteCmd = &cobra.Command{
	Use:   "delete <KEY>",
	Short: "Delete a secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		env := getEnv(cmd)
		key := args[0]
		client := newAuthClient()

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if !promptConfirm(fmt.Sprintf("Delete '%s' from %s?", key, env)) {
				info("Cancelled.")
				return
			}
		}

		sp := startSpinner(fmt.Sprintf("Deleting %s...", key))
		_, _, err := client.request("DELETE", fmt.Sprintf("/api/v1/projects/%s/secrets/%s?environment=%s", slug, key, env), nil)
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop(fmt.Sprintf("Deleted %s from %s", key, env))
	},
}

func init() {
	secretsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}

func parseKeyValue(s string) (string, string, bool) {
	idx := strings.Index(s, "=")
	if idx <= 0 {
		return "", "", false
	}
	return s[:idx], s[idx+1:], true
}

// Status command — shows current project config
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current project and environment status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		printKeyValue("API", getAPIURL())

		creds, err := loadCredentials()
		if err != nil {
			printKeyValue("Auth", red.Sprint("not logged in"))
		} else {
			printKeyValue("Auth", green.Sprintf("logged in as %s", creds.Email))
		}

		slug := ""
		func() {
			defer func() { recover() }() // getProjectSlug calls os.Exit
			slug = getProjectSlug()
		}()
		if slug != "" {
			printKeyValue("Project", slug)
		} else {
			printKeyValue("Project", dim.Sprint("not linked (run 'envault init <slug>')"))
		}
		fmt.Println()
	},
}
