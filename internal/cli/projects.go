package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Run:   projectsListRun,
}

func init() {
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)
}

func projectsListRun(cmd *cobra.Command, args []string) {
	client := newAuthClient()

	sp := startSpinner("Loading projects...")
	body, _, err := client.request("GET", "/api/v1/projects", nil)
	if err != nil {
		sp.fail("Failed: " + err.Error())
		os.Exit(1)
	}
	sp.stop("Projects loaded")
	fmt.Println()

	var projects []struct {
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		CreatedAt string `json:"created_at"`
	}
	json.Unmarshal(body, &projects)

	if len(projects) == 0 {
		info("No projects yet. Create one with 'envault projects create <name>'")
		fmt.Println()
		return
	}

	rows := make([][]string, len(projects))
	for i, p := range projects {
		rows[i] = []string{p.Name, p.Slug, p.CreatedAt[:10]}
	}
	printTable([]string{"Name", "Slug", "Created"}, rows)
	fmt.Println()
	dim.Printf("  %d project(s). Use 'envault init <slug>' to link one.\n\n", len(projects))
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := ""
		if len(args) > 0 {
			name = args[0]
		} else {
			name = prompt("Project name")
		}

		if name == "" {
			fatal("Project name is required")
		}

		client := newAuthClient()

		sp := startSpinner("Creating project...")
		body, _, err := client.request("POST", "/api/v1/projects", map[string]string{"name": name})
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Project created")
		fmt.Println()

		var resp struct {
			Project struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"project"`
			VaultToken string `json:"vault_token"`
		}
		json.Unmarshal(body, &resp)

		printKeyValue("Name", resp.Project.Name)
		printKeyValue("Slug", resp.Project.Slug)
		fmt.Println()

		if resp.VaultToken != "" {
			yellow.Println("  ⚠ Save this Vault token — it won't be shown again:")
			fmt.Println()
			boldCyan.Printf("    %s\n", resp.VaultToken)
			fmt.Println()
		}

		info("Link this project: envault init %s", resp.Project.Slug)
		fmt.Println()
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete <slug>",
	Short: "Delete a project permanently",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]

		if !promptConfirm(fmt.Sprintf("Delete project '%s' permanently?", slug)) {
			info("Cancelled.")
			return
		}

		client := newAuthClient()

		sp := startSpinner("Deleting project...")
		_, _, err := client.request("DELETE", "/api/v1/projects/"+slug, nil)
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Project deleted")
	},
}
