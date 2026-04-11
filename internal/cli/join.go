package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join <invite-code>",
	Short: "Join a project using an invite code",
	Long: `Accept a team invitation using the code shared by a project admin.

You must be logged in first (envault login or envault signup).

Example:
  envault join A1B2C3D4`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := args[0]
		client := newAuthClient()

		sp := startSpinner("Accepting invite...")
		body, _, err := client.request("POST", "/api/v1/invite/accept", map[string]string{
			"invite_code": code,
		})
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Invite accepted")
		fmt.Println()

		var resp struct {
			Member struct {
				Role string `json:"role"`
				User struct {
					Email string `json:"email"`
				} `json:"user"`
			} `json:"member"`
		}
		json.Unmarshal(body, &resp)

		printKeyValue("Role", resp.Member.Role)
		fmt.Println()
		success("You're now a team member!")
		dim.Println("  Run 'envault projects' to see your projects.")
		fmt.Println()
	},
}
