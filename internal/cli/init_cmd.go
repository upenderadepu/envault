package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initProjectCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initialize a new project",
	Long:  "Creates a new project on the Envault server and writes the configuration to ~/.envault.yaml.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		client := &cliClient{
			baseURL: getAPIURL(),
			token:   getAuthToken(),
			http:    &http.Client{},
		}

		body, _, err := client.request("POST", "/api/v1/projects", map[string]string{"name": name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var resp struct {
			Project struct {
				Slug string `json:"slug"`
			} `json:"project"`
			VaultToken string `json:"vault_token"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		// Write config
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".envault.yaml")
		configContent := fmt.Sprintf("api_url: %s\nproject_slug: %s\nvault_token: %s\n",
			getAPIURL(), resp.Project.Slug, resp.VaultToken)
		if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not write config file: %v\n", err)
		}

		fmt.Printf("Project created: %s\n", resp.Project.Slug)
		fmt.Printf("Config written to: %s\n", configPath)
		fmt.Printf("\nVault Token (save this — shown only once):\n  %s\n", resp.VaultToken)
	},
}
