package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse represents the login response
type TokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		User struct {
			ID    uint   `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"user"`
		Tokens struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
		} `json:"tokens"`
	} `json:"data"`
}

// ProfileResponse represents the profile response
type ProfileResponse struct {
	Status string `json:"status"`
	Data   struct {
		User struct {
			ID            uint    `json:"id"`
			Name          string  `json:"name"`
			Email         string  `json:"email"`
			Age           int64   `json:"age"`
			Gender        string  `json:"gender"`
			Weight        float64 `json:"weight"`
			Height        float64 `json:"height"`
			Goal          string  `json:"goal"`
			ActivityLevel string  `json:"activity_level"`
			Role          string  `json:"role"`
		} `json:"user"`
	} `json:"data"`
}

// Login authenticates a user and returns a JWT token
func Login(email, password string) (*TokenResponse, error) {
	// Create login request
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &tokenResp, nil
}

// FetchProfile retrieves the user profile using the access token
func FetchProfile(accessToken string) (*ProfileResponse, error) {
	// Create HTTP request - updated path to avoid conflicts
	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("profile fetch failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var profileResp ProfileResponse
	if err := json.Unmarshal(body, &profileResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &profileResp, nil
}

func main() {
	// Example credentials - replace with actual values
	email := "user@example.com"
	password := "password123"

	// Login to get tokens
	fmt.Println("Logging in...")
	tokenResp, err := Login(email, password)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		return
	}

	fmt.Printf("Login successful for user: %s (ID: %d)\n", 
		tokenResp.Data.User.Name, 
		tokenResp.Data.User.ID)
	
	fmt.Printf("Access token: %s\n", tokenResp.Data.Tokens.AccessToken[:20]+"...")
	fmt.Printf("Expires in: %d seconds\n", tokenResp.Data.Tokens.ExpiresIn)

	// Use the token to fetch the user profile
	fmt.Println("\nFetching user profile...")
	profileResp, err := FetchProfile(tokenResp.Data.Tokens.AccessToken)
	if err != nil {
		fmt.Printf("Profile fetch error: %v\n", err)
		return
	}

	// Display user profile
	user := profileResp.Data.User
	fmt.Printf("User Profile:\n")
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Age: %d\n", user.Age)
	fmt.Printf("  Gender: %s\n", user.Gender)
	fmt.Printf("  Weight: %.1f\n", user.Weight)
	fmt.Printf("  Height: %.1f\n", user.Height)
	fmt.Printf("  Goal: %s\n", user.Goal)
	fmt.Printf("  Activity Level: %s\n", user.ActivityLevel)
} 