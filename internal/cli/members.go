package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage team members",
	Run:   membersListRun,
}

func init() {
	membersCmd.AddCommand(membersAddCmd)
	membersCmd.AddCommand(membersRemoveCmd)
}

func membersListRun(cmd *cobra.Command, args []string) {
	slug := getProjectSlug()
	client := newAuthClient()

	sp := startSpinner("Loading members...")
	body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/members", slug), nil)
	if err != nil {
		sp.fail("Failed: " + err.Error())
		os.Exit(1)
	}
	sp.stop("Members loaded")
	fmt.Println()

	var members []struct {
		ID   string `json:"id"`
		Role string `json:"role"`
		User struct {
			Email string `json:"email"`
		} `json:"user"`
		IsActive  bool   `json:"is_active"`
		InvitedAt string `json:"invited_at"`
	}
	json.Unmarshal(body, &members)

	if len(members) == 0 {
		info("No team members. Invite someone with 'envault members add <email>'")
		fmt.Println()
		return
	}

	rows := make([][]string, len(members))
	for i, m := range members {
		email := m.User.Email
		if email == "" {
			email = "(pending)"
		}
		status := green.Sprint("active")
		if !m.IsActive {
			status = dim.Sprint("invited")
		}
		rows[i] = []string{email, m.Role, status}
	}
	printTable([]string{"Email", "Role", "Status"}, rows)
	fmt.Println()
}

var membersAddCmd = &cobra.Command{
	Use:   "add [email]",
	Short: "Invite a team member",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		client := newAuthClient()

		email := ""
		if len(args) > 0 {
			email = args[0]
		} else {
			email = prompt("Email")
		}

		role, _ := cmd.Flags().GetString("role")
		if role == "" {
			fmt.Println()
			fmt.Println("  Roles:")
			fmt.Println("    admin      - Full access (read, write, delete, manage)")
			fmt.Println("    developer  - Read & write secrets")
			fmt.Println("    ci         - Read-only access")
			fmt.Println()
			role = prompt("Role (admin/developer/ci)")
		}

		sp := startSpinner(fmt.Sprintf("Inviting %s...", email))
		body, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/members", slug), map[string]string{
			"email": email,
			"role":  role,
		})
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Member invited")
		fmt.Println()

		var resp struct {
			InviteCode string `json:"invite_code"`
		}
		json.Unmarshal(body, &resp)

		if resp.InviteCode != "" {
			fmt.Println()
			yellow.Println("  ⚠ Share this invite code with the member:")
			fmt.Println()
			boldCyan.Printf("    %s\n", resp.InviteCode)
			fmt.Println()
			dim.Println("  They can join by running:")
			cyan.Printf("    envault join %s\n", resp.InviteCode)
			fmt.Println()
			dim.Println("  This code is single-use and won't be shown again.")
		}
		fmt.Println()
	},
}

func init() {
	membersAddCmd.Flags().String("role", "", "Member role (admin, developer, ci)")
}

var membersRemoveCmd = &cobra.Command{
	Use:   "remove <email-or-id>",
	Short: "Remove a team member",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		client := newAuthClient()
		target := args[0]

		// If it looks like an email, find the member ID first
		memberID := target
		if contains(target, "@") {
			body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/members", slug), nil)
			if err != nil {
				fatal("Failed to list members: %s", err.Error())
			}
			var members []struct {
				ID   string `json:"id"`
				User struct {
					Email string `json:"email"`
				} `json:"user"`
			}
			json.Unmarshal(body, &members)
			for _, m := range members {
				if m.User.Email == target {
					memberID = m.ID
					break
				}
			}
			if memberID == target {
				fatal("Member '%s' not found", target)
			}
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if !promptConfirm(fmt.Sprintf("Remove '%s'? Their access will be revoked.", target)) {
				info("Cancelled.")
				return
			}
		}

		sp := startSpinner("Removing member...")
		_, _, err := client.request("DELETE", fmt.Sprintf("/api/v1/projects/%s/members/%s", slug, memberID), nil)
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Member removed")
	},
}

func init() {
	membersRemoveCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
