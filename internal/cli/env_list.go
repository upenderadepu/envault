package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables",
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets in an environment",
	Long:  "Lists all secret keys (without values) for the given environment.",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		slug := getProjectSlug()
		client := newCLIClient()

		body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/secrets?environment=%s", slug, env), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var secrets []struct {
			KeyName        string `json:"key_name"`
			VaultVersion   int    `json:"vault_version"`
			LastModifiedAt string `json:"last_modified_at"`
		}
		if err := json.Unmarshal(body, &secrets); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		if len(secrets) == 0 {
			fmt.Printf("No secrets in %s/%s\n", slug, env)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tVERSION\tMODIFIED")
		fmt.Fprintln(w, "---\t-------\t--------")
		for _, s := range secrets {
			fmt.Fprintf(w, "%s\tv%d\t%s\n", s.KeyName, s.VaultVersion, s.LastModifiedAt)
		}
		w.Flush()

		fmt.Printf("\n%d secret(s) in %s/%s\n", len(secrets), slug, env)
	},
}

func init() {
	envListCmd.Flags().String("env", "development", "Environment (development, staging, production)")
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envPullCmd)
	envCmd.AddCommand(envPushCmd)
}
