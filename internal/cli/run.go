package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"syscall"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [--env environment] -- <command> [args...]",
	Short: "Run a command with secrets injected as environment variables",
	Long: `Fetches all secrets for the given environment and injects them
as environment variables into the specified command.

Examples:
  envault run -- npm start
  envault run --env production -- node server.js
  envault run -- docker compose up`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse our own flags before --
		env := "development"
		commandArgs := args

		for i := 0; i < len(args); i++ {
			if args[i] == "--" {
				commandArgs = args[i+1:]
				break
			}
			if args[i] == "--env" && i+1 < len(args) {
				env = args[i+1]
				i++
				continue
			}
		}

		if len(commandArgs) == 0 {
			fatal("No command specified. Usage: envault run [--env ENV] -- <command>")
		}

		slug := getProjectSlug()
		client := newAuthClient()

		// Fetch secret keys
		sp := startSpinner(fmt.Sprintf("Fetching %s secrets...", env))
		body, _, err := client.request("GET", fmt.Sprintf("/api/v1/projects/%s/secrets?environment=%s", slug, env), nil)
		if err != nil {
			sp.fail("Failed to fetch secrets: " + err.Error())
			os.Exit(1)
		}

		var secrets []struct {
			KeyName string `json:"key_name"`
		}
		json.Unmarshal(body, &secrets)

		// Fetch each value
		envVars := make(map[string]string)
		for _, s := range secrets {
			valBody, _, err := client.request("GET",
				fmt.Sprintf("/api/v1/projects/%s/secrets/%s?environment=%s", slug, s.KeyName, env), nil)
			if err != nil {
				sp.fail(fmt.Sprintf("Failed to fetch %s: %s", s.KeyName, err.Error()))
				os.Exit(1)
			}
			var val struct {
				Value string `json:"value"`
			}
			json.Unmarshal(valBody, &val)
			envVars[s.KeyName] = val.Value
		}
		sp.stop(fmt.Sprintf("Loaded %d secrets from %s", len(envVars), env))

		// Print injected keys
		if len(envVars) > 0 {
			keys := make([]string, 0, len(envVars))
			for k := range envVars {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			dim.Printf("  Injecting: %s\n", join(keys))
		}

		fmt.Println()
		dim.Printf("  $ %s\n\n", joinArgs(commandArgs))

		// Build environment
		environ := os.Environ()
		for k, v := range envVars {
			environ = append(environ, fmt.Sprintf("%s=%s", k, v))
		}

		// Execute command
		binary, err := exec.LookPath(commandArgs[0])
		if err != nil {
			fatal("Command not found: %s", commandArgs[0])
		}

		process := &exec.Cmd{
			Path:   binary,
			Args:   commandArgs,
			Env:    environ,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		// Forward signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-sigCh
			if process.Process != nil {
				process.Process.Signal(sig)
			}
		}()

		if err := process.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fatal("Command failed: %v", err)
		}
	},
}

func join(ss []string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}
