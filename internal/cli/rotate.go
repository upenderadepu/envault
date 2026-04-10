package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate project credentials",
	Long:  "Rotates Vault credentials for the current project. All existing tokens will be invalidated.",
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		client := newCLIClient()

		body, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/rotate", slug), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var resp struct {
			NewToken string `json:"new_vault_token"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Credentials rotated successfully.")
		fmt.Printf("\nNew Vault Token (save this — shown only once):\n  %s\n", resp.NewToken)
		fmt.Println("\nUpdate your ~/.envault.yaml with the new token.")
	},
}
