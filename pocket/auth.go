package pocket

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type loginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Record struct {
		Avatar          string `json:"avatar"`
		CollectionID    string `json:"collectionId"`
		CollectionName  string `json:"collectionName"`
		Created         string `json:"created"`
		Email           string `json:"email"`
		EmailVisibility bool   `json:"emailVisibility"`
		ID              string `json:"id"`
		Name            string `json:"name"`
		Updated         string `json:"updated"`
		Username        string `json:"username"`
		Verified        bool   `json:"verified"`
	} `json:"record"`
	Token string `json:"token"`
}

type User struct {
	ID        string
	Name      string
	Email     string
	Avatar    string
	Token     string
	CreatedAt time.Time
}

func loginPrompt() ([]byte, error) {
	var input loginRequest
	prompt := []*survey.Question{
		{
			Name:   "identity",
			Prompt: &survey.Input{Message: "Identity:"},
		},
		{
			Name:   "password",
			Prompt: &survey.Password{Message: "Password:"},
		},
	}
	if err := survey.Ask(prompt, &input); err != nil {
		return nil, fmt.Errorf("failed to prompt for credentials: %v", err)
	}

	// Create a new login request
	reqBody := loginRequest{
		Identity: input.Identity,
		Password: input.Password,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode login request: %v", err)
	}

	return reqBytes, nil
}

func getAuthToken(loginRequest []byte) (*AuthResponse, error) {
	domain := os.Getenv("DOMAIN")
	url := fmt.Sprintf("https://%s/api/collections/users/auth-with-password", domain)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginRequest))
	if err != nil {
		return nil, fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authenticate user: %v", resp.Status)
	}

	// Decode the login response
	var loginResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %v", err)
	}

	return &loginResp, nil
}

func saveUser(authResp *AuthResponse) error {
	// Create a new `User` struct with the authenticated user's details
	currentUser := &User{
		ID:        authResp.Record.ID,
		Name:      authResp.Record.Name,
		Email:     authResp.Record.Email,
		Avatar:    authResp.Record.Avatar,
		Token:     authResp.Token,
		CreatedAt: time.Now(),
	}
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath()), 0755); err != nil {
		return fmt.Errorf("failed to create directory for user file: %v", err)
	}
	// Create a new file for the user's details
	f, err := os.Create(configPath())
	if err != nil {
		return fmt.Errorf("failed to create user file: %v", err)
	}
	defer f.Close()

	// Encode the user's details and write them to the file
	enc := gob.NewEncoder(f)
	if err := enc.Encode(currentUser); err != nil {
		return fmt.Errorf("failed to encode user details: %v", err)
	}

	return nil
}

func LoadUser() (*User, error) {
	// Open the file containing the user's details
	f, err := os.Open(configPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open user file: %v", err)
	}
	defer f.Close()

	// Decode the user's details from the file
	dec := gob.NewDecoder(f)
	var user User
	if err := dec.Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user details: %v", err)
	}

	return &user, nil
}

func configPath() string {
	osCurrentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := osCurrentUser.HomeDir
	userFile := fmt.Sprintf("%s/.config/go-journal/user", homeDir)
	return userFile
}

func Auth_refresh() error {
	user, _ := LoadUser()
	domain := os.Getenv("DOMAIN")
	url := fmt.Sprintf("https://%s/api/collections/users/auth-refresh", domain)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %v", err)
	}

	// Add the authentication token to the request headers
	req.Header.Set("Authorization", user.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refresh request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh authentication token: %v", resp.Status)
	}

	return nil
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Pocket",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a prompt for the identity and password
		loginReq, err := loginPrompt()
		if err != nil {
			return fmt.Errorf("failed to prompt for login credentials: %v", err)
		}
		// Send a POST request to the login endpoint
		authResp, err := getAuthToken(loginReq)
		if err != nil {
			return fmt.Errorf("failed to authenticate user: %v", err)
		}

		// Save the authenticated user's details for future use
		if err := saveUser(authResp); err != nil {
			return fmt.Errorf("failed to save user details: %v", err)
		}
		fmt.Printf("Logged in as %s (%s)\n", authResp.Record.Name, authResp.Record.Email)
		return nil
	},
}
