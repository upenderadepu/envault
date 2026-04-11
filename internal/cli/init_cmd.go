package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initProjectCmd = &cobra.Command{
	Use:   "init <project-slug>",
	Short: "Link current directory to a project",
	Long: `Creates a .envault.yaml config in the current directory.
Get the project slug from 'envault projects' or the dashboard.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]

		// Verify the project exists
		client := newAuthClient()
		sp := startSpinner("Verifying project...")
		_, _, err := client.request("GET", "/api/v1/projects/"+slug, nil)
		if err != nil {
			sp.fail("Project not found: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Project verified")

		// Write config
		apiURL := getAPIURL()
		config := fmt.Sprintf("api_url: %s\nproject_slug: %s\n", apiURL, slug)
		if err := os.WriteFile(".envault.yaml", []byte(config), 0644); err != nil {
			fatal("Failed to write config: %v", err)
		}

		fmt.Println()
		success("Linked to project '%s'", slug)
		dim.Println("  Config saved to .envault.yaml")
		fmt.Println()
		fmt.Println("  Next steps:")
		cyan.Println("    envault secrets                      # List secrets")
		cyan.Println("    envault secrets set KEY=VALUE         # Set a secret")
		cyan.Println("    envault env pull                     # Pull to .env file")
		cyan.Println("    envault run -- npm start             # Run with secrets")
		fmt.Println()
	},
}
