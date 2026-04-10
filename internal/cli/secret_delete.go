package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var secretDeleteCmd = &cobra.Command{
	Use:   "delete KEY",
	Short: "Delete a secret",
	Long:  "Permanently removes a secret from the given environment.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		slug := getProjectSlug()
		client := newCLIClient()

		_, _, err := client.request("DELETE",
			fmt.Sprintf("/api/v1/projects/%s/secrets/%s?environment=%s", slug, args[0], env), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Secret %s deleted from %s/%s\n", args[0], slug, env)
	},
}

func init() {
	secretDeleteCmd.Flags().String("env", "development", "Environment (development, staging, production)")
	secretCmd.AddCommand(secretDeleteCmd)
}
