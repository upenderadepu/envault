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
	Long:  "Envault stores secrets in HashiCorp Vault, provides a CLI to inject them into any environment, and gives teams audited, role-scoped access.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("api-url", "http://localhost:8080", "Envault API server URL")
	rootCmd.PersistentFlags().String("vault-addr", "", "Vault server address")
	rootCmd.PersistentFlags().String("vault-token", "", "Vault token for authentication")
	rootCmd.PersistentFlags().String("project", "", "Project slug")

	viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("vault_addr", rootCmd.PersistentFlags().Lookup("vault-addr"))
	viper.BindPFlag("vault_token", rootCmd.PersistentFlags().Lookup("vault-token"))
	viper.BindPFlag("project_slug", rootCmd.PersistentFlags().Lookup("project"))

	// Subcommands
	rootCmd.AddCommand(initProjectCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(secretCmd)
	rootCmd.AddCommand(onboardCmd)
	rootCmd.AddCommand(rotateCmd)
}

func initConfig() {
	viper.SetConfigName(".envault")
	viper.SetConfigType("yaml")
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("ENVAULT")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config:", viper.ConfigFileUsed())
	}
}

func getAPIURL() string {
	return viper.GetString("api_url")
}

func getProjectSlug() string {
	slug := viper.GetString("project_slug")
	if slug == "" {
		fmt.Fprintln(os.Stderr, "Error: project slug is required. Use --project flag or set it in ~/.envault.yaml")
		os.Exit(1)
	}
	return slug
}

func getAuthToken() string {
	token := viper.GetString("vault_token")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: auth token is required. Use --vault-token flag or set it in ~/.envault.yaml")
		os.Exit(1)
	}
	return token
}
