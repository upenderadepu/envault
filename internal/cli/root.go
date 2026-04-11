package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "envault",
	Short: "Envault — Secure secrets management for teams",
	Long:  "Store secrets in HashiCorp Vault, inject them into any environment,\nand give teams audited, role-scoped access.",
	Run: func(cmd *cobra.Command, args []string) {
		printLogo()
		fmt.Println("  Usage:")
		fmt.Println("    envault login          Sign in to your account")
		fmt.Println("    envault projects       List all projects")
		fmt.Println("    envault secrets        List secrets in current project")
		fmt.Println("    envault run -- <cmd>   Run command with injected secrets")
		fmt.Println()
		dim.Println("  Run 'envault <command> --help' for more info.")
		fmt.Println()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("api-url", "", "Envault API server URL")
	rootCmd.PersistentFlags().String("project", "", "Project slug")
	rootCmd.PersistentFlags().String("env", "development", "Environment (development, staging, production)")

	viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("project_slug", rootCmd.PersistentFlags().Lookup("project"))

	// Commands
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(signupCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(initProjectCmd)
	rootCmd.AddCommand(secretsCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(membersCmd)
	rootCmd.AddCommand(joinCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(statusCmd)
}

func initConfig() {
	// Project-level config (.envault.yaml in current dir)
	viper.SetConfigName(".envault")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("ENVAULT")
	viper.AutomaticEnv()

	viper.ReadInConfig() // ignore error — config file is optional
}

func getAPIURL() string {
	url := viper.GetString("api_url")
	if url == "" {
		url = os.Getenv("ENVAULT_API_URL")
	}
	if url == "" {
		url = "http://localhost:8080"
	}
	return url
}

func getProjectSlug() string {
	slug := viper.GetString("project_slug")
	if slug == "" {
		fatal("No project linked. Run 'envault init <slug>' or use --project flag.")
	}
	return slug
}

func getEnv(cmd *cobra.Command) string {
	env, _ := cmd.Flags().GetString("env")
	if env == "" {
		env = viper.GetString("environment")
	}
	if env == "" {
		env = "development"
	}
	return env
}
