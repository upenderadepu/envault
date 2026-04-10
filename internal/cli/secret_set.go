package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var secretSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE",
	Short: "Set a secret",
	Long:  "Creates or updates a single secret in the given environment.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		slug := getProjectSlug()
		client := newCLIClient()

		idx := strings.Index(args[0], "=")
		if idx == -1 {
			fmt.Fprintln(os.Stderr, "Error: argument must be in KEY=VALUE format")
			os.Exit(1)
		}

		key := args[0][:idx]
		value := args[0][idx+1:]

		if key == "" {
			fmt.Fprintln(os.Stderr, "Error: key cannot be empty")
			os.Exit(1)
		}

		payload := map[string]string{
			"environment": env,
			"key":         key,
			"value":       value,
		}

		body, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/secrets", slug), payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var meta struct {
			KeyName      string `json:"key_name"`
			VaultVersion int    `json:"vault_version"`
		}
		json.Unmarshal(body, &meta)

		fmt.Printf("Secret %s set (version %d) in %s/%s\n", meta.KeyName, meta.VaultVersion, slug, env)
	},
}

func init() {
	secretSetCmd.Flags().String("env", "development", "Environment (development, staging, production)")
	secretCmd.AddCommand(secretSetCmd)
}
