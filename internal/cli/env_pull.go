package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var envPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull secrets to a .env file",
	Long:  "Downloads all secrets for the given environment and writes them to a .env file.",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		output, _ := cmd.Flags().GetString("output")
		slug := getProjectSlug()
		client := newCLIClient()

		// List secrets
		body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/secrets?environment=%s", slug, env), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var secrets []struct {
			KeyName string `json:"key_name"`
		}
		if err := json.Unmarshal(body, &secrets); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		if len(secrets) == 0 {
			fmt.Println("No secrets found.")
			return
		}

		// Fetch each value
		var lines []string
		for _, s := range secrets {
			valBody, _, err := client.request("GET",
				fmt.Sprintf("/api/v1/projects/%s/secrets/%s?environment=%s", slug, s.KeyName, env), nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not fetch %s: %v\n", s.KeyName, err)
				continue
			}
			var val struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(valBody, &val); err != nil {
				continue
			}
			lines = append(lines, fmt.Sprintf("%s=%s", s.KeyName, val.Value))
		}

		sort.Strings(lines)

		content := ""
		for _, l := range lines {
			content += l + "\n"
		}

		if err := os.WriteFile(output, []byte(content), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Pulled %d secrets to %s\n", len(lines), output)
	},
}

func init() {
	envPullCmd.Flags().String("env", "development", "Environment (development, staging, production)")
	envPullCmd.Flags().StringP("output", "o", ".env", "Output file path")
}
