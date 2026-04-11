package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Email        string `json:"email"`
	APIURL       string `json:"api_url"`
	SupabaseURL  string `json:"supabase_url"`
}

func credentialsPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".envault")
	os.MkdirAll(dir, 0700)
	return filepath.Join(dir, "credentials.json")
}

func saveCredentials(creds *credentials) error {
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsPath(), data, 0600)
}

func loadCredentials() (*credentials, error) {
	data, err := os.ReadFile(credentialsPath())
	if err != nil {
		return nil, err
	}
	var creds credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

func deleteCredentials() {
	os.Remove(credentialsPath())
}

// getAuthToken returns a valid JWT token, refreshing if needed.
func getAuthToken() string {
	creds, err := loadCredentials()
	if err != nil {
		fatal("Not logged in. Run 'envault login' first.")
	}

	// Check if token is expired (with 60s buffer)
	if time.Now().Unix() > creds.ExpiresAt-60 {
		// Try to refresh
		newCreds, err := refreshToken(creds)
		if err != nil {
			fatal("Session expired. Run 'envault login' to re-authenticate.")
		}
		creds = newCreds
	}

	return creds.AccessToken
}

func supabaseAuth(supabaseURL, email, password string, isSignUp bool) (*credentials, error) {
	endpoint := "/auth/v1/token?grant_type=password"
	if isSignUp {
		endpoint = "/auth/v1/signup"
	}

	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	req, _ := http.NewRequest("POST", supabaseURL+endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", getSupabaseAnonKey())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errResp struct {
			Msg          string `json:"msg"`
			ErrorMessage string `json:"error_description"`
		}
		json.Unmarshal(respBody, &errResp)
		msg := errResp.Msg
		if msg == "" {
			msg = errResp.ErrorMessage
		}
		if msg == "" {
			msg = "authentication failed"
		}
		return nil, fmt.Errorf("%s", msg)
	}

	var authResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		User         struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, fmt.Errorf("invalid response")
	}

	return &credentials{
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + authResp.ExpiresIn,
		Email:        authResp.User.Email,
		APIURL:       getAPIURL(),
		SupabaseURL:  supabaseURL,
	}, nil
}

func refreshToken(creds *credentials) (*credentials, error) {
	body, _ := json.Marshal(map[string]string{
		"refresh_token": creds.RefreshToken,
	})

	url := creds.SupabaseURL + "/auth/v1/token?grant_type=refresh_token"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", getSupabaseAnonKey())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("refresh failed")
	}

	var authResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &authResp)

	newCreds := &credentials{
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + authResp.ExpiresIn,
		Email:        creds.Email,
		APIURL:       creds.APIURL,
		SupabaseURL:  creds.SupabaseURL,
	}
	saveCredentials(newCreds)
	return newCreds, nil
}

func getSupabaseURL() string {
	url := viper.GetString("supabase_url")
	if url == "" {
		url = os.Getenv("ENVAULT_SUPABASE_URL")
	}
	if url == "" {
		// Try loading from saved credentials
		if creds, err := loadCredentials(); err == nil && creds.SupabaseURL != "" {
			return creds.SupabaseURL
		}
	}
	return url
}

func getSupabaseAnonKey() string {
	key := viper.GetString("supabase_anon_key")
	if key == "" {
		key = os.Getenv("ENVAULT_SUPABASE_ANON_KEY")
	}
	return key
}

// Login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in to your Envault account",
	Run: func(cmd *cobra.Command, args []string) {
		printLogo()

		supabaseURL := getSupabaseURL()
		if supabaseURL == "" {
			supabaseURL = prompt("Supabase URL (e.g. https://xxx.supabase.co)")
		}
		if getSupabaseAnonKey() == "" {
			fatal("ENVAULT_SUPABASE_ANON_KEY not set. Export it or add supabase_anon_key to .envault.yaml")
		}

		email := prompt("Email")
		password := promptPassword("Password")

		sp := startSpinner("Signing in...")
		creds, err := supabaseAuth(supabaseURL, email, password, false)
		if err != nil {
			sp.fail("Login failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Logged in")

		if err := saveCredentials(creds); err != nil {
			fatal("Failed to save credentials: %v", err)
		}

		fmt.Println()
		success("Welcome back, %s!", creds.Email)
		dim.Println("  Credentials saved to ~/.envault/credentials.json")
		fmt.Println()
	},
}

// Signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new Envault account",
	Run: func(cmd *cobra.Command, args []string) {
		printLogo()

		supabaseURL := getSupabaseURL()
		if supabaseURL == "" {
			supabaseURL = prompt("Supabase URL (e.g. https://xxx.supabase.co)")
		}
		if getSupabaseAnonKey() == "" {
			fatal("ENVAULT_SUPABASE_ANON_KEY not set. Export it or add supabase_anon_key to .envault.yaml")
		}

		email := prompt("Email")
		password := promptPassword("Password")

		sp := startSpinner("Creating account...")
		creds, err := supabaseAuth(supabaseURL, email, password, true)
		if err != nil {
			sp.fail("Signup failed: " + err.Error())
			os.Exit(1)
		}
		sp.stop("Account created")

		if creds.AccessToken != "" {
			saveCredentials(creds)
			success("You're all set, %s!", creds.Email)
		} else {
			info("Check your email to confirm your account, then run 'envault login'.")
		}
		fmt.Println()
	},
}

// Logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out and clear stored credentials",
	Run: func(cmd *cobra.Command, args []string) {
		deleteCredentials()
		success("Logged out.")
	},
}

// Whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current authenticated user",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := loadCredentials()
		if err != nil {
			fatal("Not logged in. Run 'envault login' first.")
		}
		fmt.Println()
		printKeyValue("Email", creds.Email)
		printKeyValue("API", creds.APIURL)

		if time.Now().Unix() > creds.ExpiresAt {
			warn("Session expired. Run 'envault login' to refresh.")
		} else {
			remaining := time.Until(time.Unix(creds.ExpiresAt, 0)).Round(time.Minute)
			printKeyValue("Session", fmt.Sprintf("valid for %s", remaining))
		}
		fmt.Println()
	},
}
