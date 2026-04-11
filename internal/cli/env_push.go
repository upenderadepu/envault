package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var envPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a .env file to the server",
	Run: func(cmd *cobra.Command, args []string) {
		slug := getProjectSlug()
		env, _ := cmd.Flags().GetString("env")
		input, _ := cmd.Flags().GetString("input")
		client := newAuthClient()

		data, err := os.ReadFile(input)
		if err != nil {
			fatal("Failed to read %s: %v", input, err)
		}

		secrets := parseEnvFile(string(data))
		if len(secrets) == 0 {
			fatal("No valid KEY=VALUE pairs found in %s", input)
		}

		if !promptConfirm(fmt.Sprintf("Push %d secrets to %s?", len(secrets), env)) {
			info("Cancelled.")
			return
		}

		sp := startSpinner(fmt.Sprintf("Pushing %d secrets to %s...", len(secrets), env))
		_, _, err = client.request("POST", fmt.Sprintf("/api/v1/projects/%s/secrets/bulk", slug), map[string]interface{}{
			"environment": env,
			"secrets":     secrets,
		})
		if err != nil {
			sp.fail("Failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop(fmt.Sprintf("Pushed %d secrets to %s", len(secrets), env))
	},
}

func init() {
	envPushCmd.Flags().String("env", "development", "Environment")
	envPushCmd.Flags().StringP("input", "i", ".env", "Input file path")
}

func parseEnvFile(content string) map[string]string {
	secrets := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		// Remove surrounding quotes
		if len(value) >= 2 &&
			((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		secrets[key] = value
	}
	return secrets
}
