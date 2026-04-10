package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var secretGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get a secret value",
	Long:  "Retrieves and displays the value of a single secret.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		slug := getProjectSlug()
		client := newCLIClient()

		body, _, err := client.request("GET",
			fmt.Sprintf("/api/v1/projects/%s/secrets/%s?environment=%s", slug, args[0], env), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var secret struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Version int    `json:"version"`
		}
		if err := json.Unmarshal(body, &secret); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		valueOnly, _ := cmd.Flags().GetBool("value-only")
		if valueOnly {
			fmt.Print(secret.Value)
		} else {
			fmt.Printf("%s=%s\n", args[0], secret.Value)
		}
	},
}

func init() {
	secretGetCmd.Flags().String("env", "development", "Environment (development, staging, production)")
	secretGetCmd.Flags().Bool("value-only", false, "Print only the value (useful for scripts)")
	secretCmd.AddCommand(secretGetCmd)
}
