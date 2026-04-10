package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var envPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a .env file to the server",
	Long:  "Reads a .env file and bulk-uploads all key-value pairs to the given environment.",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		file, _ := cmd.Flags().GetString("file")
		slug := getProjectSlug()
		client := newCLIClient()

		secrets, err := parseEnvFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(secrets) == 0 {
			fmt.Println("No secrets found in file.")
			return
		}

		payload := map[string]interface{}{
			"environment": env,
			"secrets":     secrets,
		}

		body, _, err := client.request("POST", fmt.Sprintf("/api/v1/projects/%s/secrets/bulk", slug), payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var result []json.RawMessage
		json.Unmarshal(body, &result)

		fmt.Printf("Pushed %d secrets to %s/%s\n", len(result), slug, env)
	},
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", path, err)
	}
	defer f.Close()

	secrets := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		// Strip surrounding quotes
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		if key != "" {
			secrets[key] = value
		}
	}

	return secrets, scanner.Err()
}

func init() {
	envPushCmd.Flags().String("env", "development", "Environment (development, staging, production)")
	envPushCmd.Flags().StringP("file", "f", ".env", "Path to .env file")
}
