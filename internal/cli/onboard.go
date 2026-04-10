package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard <email>",
	Short: "Add a team member to the project",
	Long:  "Invites a user by email and assigns them a role in the current project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		role, _ := cmd.Flags().GetString("role")
		slug := getProjectSlug()
		client := newCLIClient()

		payload := map[string]string{
			"email": args[0],
			"role":  role,
		}

		body, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/members", slug), payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var resp struct {
			Member struct {
				Email string `json:"email"`
				Role  string `json:"role"`
			} `json:"member"`
			VaultToken string `json:"vault_token"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added %s as %s to %s\n", resp.Member.Email, resp.Member.Role, slug)
		fmt.Printf("\nVault Token (share securely — shown only once):\n  %s\n", resp.VaultToken)
	},
}

func init() {
	onboardCmd.Flags().String("role", "developer", "Role (admin, developer, ci)")
}
