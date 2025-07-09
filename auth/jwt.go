package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService provides methods for JWT token operations
type JWTService struct {
	config Config
}

// NewJWTService creates a new JWTService instance
func NewJWTService() *JWTService {
	return &JWTService{
		config: GetConfig(),
	}
}

// TokenPair represents a pair of JWT tokens (access and refresh)
type TokenPair struct {
	AccessTokenID  string `json:"access_token_id"`
	RefreshTokenID string `json:"refresh_token_id"`
	ExpiresIn      int64  `json:"expires_in"` // Seconds until access token expires
}

// GenerateTokenPair creates a new access and refresh token pair
func (s *JWTService) GenerateTokenPair(userID uint, email, role string) (TokenPair, error) {
	// Create access token
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(s.config.TokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.config.Issuer,
		"type":    "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return TokenPair{}, err
	}

	// Create refresh token with longer expiry
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.RefreshExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.config.Issuer,
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return TokenPair{}, err
	}

	// Store tokens in Redis and get their IDs
	accessTokenID, err := StoreToken(accessTokenString, s.config.TokenExpiry)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to store access token: %v", err)
	}

	refreshTokenID, err := StoreToken(refreshTokenString, s.config.RefreshExpiry)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to store refresh token: %v", err)
	}

	// Calculate seconds until access token expiry
	expiresIn := int64(s.config.TokenExpiry / time.Second)

	return TokenPair{
		AccessTokenID:  accessTokenID,
		RefreshTokenID: refreshTokenID,
		ExpiresIn:      expiresIn,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenID string) (*jwt.Token, jwt.MapClaims, error) {
	// Get token from Redis
	tokenString, err := GetToken(tokenID)
	if err != nil {
		return nil, nil, fmt.Errorf("token not found: %v", err)
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Explicit expiration check
		if exp, ok := claims["exp"].(float64); ok {
			currentTime := time.Now().Unix()
			expirationTime := int64(exp)

			if currentTime > expirationTime {
				// Token has expired, also remove it from Redis to clean up
				err := DeleteToken(tokenID)
				if err != nil {
					return nil, nil, err
				}
				return nil, nil, fmt.Errorf("token has expired at %v (current time: %v)",
					time.Unix(expirationTime, 0), time.Unix(currentTime, 0))
			}
		} else {
			return nil, nil, errors.New("token missing expiration claim")
		}

		// Validate issued at time (optional security check)
		if iat, ok := claims["iat"].(float64); ok {
			issuedTime := int64(iat)
			currentTime := time.Now().Unix()

			// Reject tokens issued in the future (clock skew tolerance of 5 minutes)
			if issuedTime > currentTime+300 {
				return nil, nil, errors.New("token issued in the future")
			}
		}

		if token.Valid {
			return token, claims, nil
		}
	}

	return nil, nil, errors.New("invalid token")
}

// ExtractClaims extracts user claims from a validated token
func (s *JWTService) ExtractClaims(claims jwt.MapClaims) (Claims, error) {
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return Claims{}, errors.New("invalid user_id in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return Claims{}, errors.New("invalid email in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return Claims{}, errors.New("invalid role in token")
	}

	return Claims{
		UserID: uint(userID),
		Email:  email,
		Role:   role,
	}, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *JWTService) RefreshAccessToken(refreshTokenID string) (TokenPair, error) {
	// Validate the refresh token
	token, claims, err := s.ValidateToken(refreshTokenID)
	if err != nil {
		return TokenPair{}, err
	}

	// Ensure token is valid
	if !token.Valid {
		return TokenPair{}, errors.New("invalid refresh token")
	}

	// Ensure this is a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return TokenPair{}, errors.New("invalid token type")
	}

	// Extract the user ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return TokenPair{}, errors.New("invalid token claims")
	}

	// Generate a new token pair
	return s.GenerateTokenPair(uint(userID), "", "")
}
