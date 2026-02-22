package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	config Config
}

// TokenPair holds signed JWT strings returned after successful login
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

func NewJWTService() *JWTService {
	return &JWTService{config: GetConfig()}
}

// GenerateTokenPair creates a signed access + refresh token pair.
// Both tokens are returned as strings and stored in HttpOnly cookies by the caller —
// no server-side token store (Redis) is needed.
func (s *JWTService) GenerateTokenPair(userID uint, email, role string) (TokenPair, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     now.Add(s.config.TokenExpiry).Unix(),
		"iat":     now.Unix(),
		"iss":     s.config.Issuer,
		"type":    "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token includes email and role so they are available when re-issuing
	// an access token without an extra database round-trip.
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     now.Add(s.config.RefreshExpiry).Unix(),
		"iat":     now.Unix(),
		"iss":     s.config.Issuer,
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.config.TokenExpiry / time.Second),
	}, nil
}

// ValidateToken parses and validates a JWT string, returning the token and its claims.
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}

// ExtractClaims extracts typed user identity from validated MapClaims.
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

// RefreshAccessToken validates a refresh token string and issues a new token pair.
// Email and role are read directly from the refresh token claims, avoiding a DB lookup.
func (s *JWTService) RefreshAccessToken(refreshTokenString string) (TokenPair, error) {
	_, claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return TokenPair{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return TokenPair{}, errors.New("token is not a refresh token")
	}

	userClaims, err := s.ExtractClaims(claims)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to extract claims from refresh token: %w", err)
	}

	return s.GenerateTokenPair(userClaims.UserID, userClaims.Email, userClaims.Role)
}
