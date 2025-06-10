package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UserLoginRequest represents the user login credentials
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginResponse represents the user login response
type UserLoginResponse struct {
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
			AccessTokenID  string `json:"access_token_id"`
			RefreshTokenID string `json:"refresh_token_id"`
			ExpiresIn      int64  `json:"expires_in"`
		} `json:"tokens"`
	} `json:"data"`
}

// UserLogoutResponse represents the user logout response
type UserLogoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UserLogin authenticates a user using the User Module Login endpoint
func UserLogin(email, password string) (*UserLoginResponse, error) {
	// Create login request
	loginReq := UserLoginRequest{
		Email:    email,
		Password: password,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(jsonData))
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
	var loginResp UserLoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &loginResp, nil
}

// UserLogout logs out a user using the User Module Logout endpoint
func UserLogout(accessTokenID string) (*UserLogoutResponse, error) {
	// Create HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/logout", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header with token ID
	req.Header.Set("Authorization", "Bearer "+accessTokenID)
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
		return nil, fmt.Errorf("logout failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var logoutResp UserLogoutResponse
	if err := json.Unmarshal(body, &logoutResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &logoutResp, nil
}

// Example usage demonstrating the complete login/logout flow
func main() {
	email := "test@example.com"
	password := "password123"

	fmt.Printf("=== User Module Login/Logout Example ===\n\n")

	// Step 1: Login
	fmt.Printf("1. Logging in with email: %s\n", email)
	loginResp, err := UserLogin(email, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Login successful!\n")
	fmt.Printf("   User ID: %d\n", loginResp.Data.User.ID)
	fmt.Printf("   User Name: %s\n", loginResp.Data.User.Name)
	fmt.Printf("   Access Token ID: %s\n", loginResp.Data.Tokens.AccessTokenID)
	fmt.Printf("   Expires In: %d seconds\n\n", loginResp.Data.Tokens.ExpiresIn)

	// Step 2: Use the token for authenticated requests
	accessTokenID := loginResp.Data.Tokens.AccessTokenID
	fmt.Printf("2. Token ID to use for authenticated requests: %s\n\n", accessTokenID)

	// Step 3: Logout
	fmt.Printf("3. Logging out...\n")
	logoutResp, err := UserLogout(accessTokenID)
	if err != nil {
		fmt.Printf("Logout failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Logout successful!\n")
	fmt.Printf("   Message: %s\n", logoutResp.Message)

	fmt.Printf("\n=== Example completed successfully! ===\n")
} 