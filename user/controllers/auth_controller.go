package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/momokapoolz/caloriesapp/auth"
	"github.com/momokapoolz/caloriesapp/dto"
	"github.com/momokapoolz/caloriesapp/user/models"
	"github.com/momokapoolz/caloriesapp/user/repository"
	"github.com/momokapoolz/caloriesapp/user/utils"
)

// UserAuthController handles authentication endpoints (login, register, logout, refresh)
type UserAuthController struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
	config     auth.Config
}

// NewUserAuthController creates a new UserAuthController
func NewUserAuthController() *UserAuthController {
	return &UserAuthController{
		userRepo:   repository.NewUserRepository(),
		jwtService: auth.NewJWTService(),
		config:     auth.GetConfig(),
	}
}

// setTokenCookies writes the access and refresh JWTs as HttpOnly cookies.
// SameSite=Strict prevents CSRF; Secure is driven by COOKIE_SECURE env var.
func (c *UserAuthController) setTokenCookies(ctx *gin.Context, tokenPair auth.TokenPair) {
	accessMaxAge := int(tokenPair.ExpiresIn)
	refreshMaxAge := int(c.config.RefreshExpiry / time.Second)

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(auth.AccessTokenCookie, tokenPair.AccessToken,
		accessMaxAge, "/", c.config.CookieDomain, c.config.CookieSecure, true)
	ctx.SetCookie(auth.RefreshTokenCookie, tokenPair.RefreshToken,
		refreshMaxAge, "/", c.config.CookieDomain, c.config.CookieSecure, true)
}

// clearTokenCookies invalidates both auth cookies by setting MaxAge=-1
func (c *UserAuthController) clearTokenCookies(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(auth.AccessTokenCookie, "", -1, "/", c.config.CookieDomain, c.config.CookieSecure, true)
	ctx.SetCookie(auth.RefreshTokenCookie, "", -1, "/", c.config.CookieDomain, c.config.CookieSecure, true)
}

// toAuthUserResponse converts a User model to the DTO used in auth responses
func toAuthUserResponse(user *models.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Age:           user.Age,
		Gender:        user.Gender,
		Weight:        user.Weight,
		Height:        user.Height,
		Goal:          user.Goal,
		ActivityLevel: user.ActivityLevel,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
	}
}

// Register creates a new user account.
// Does NOT auto-login — the client should call POST /login after registration.
func (c *UserAuthController) Register(ctx *gin.Context) {
	var req dto.RegisterDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[Register] Invalid request: %v", err)
		ctx.JSON(http.StatusBadRequest, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Invalid request format",
		})
		return
	}

	log.Printf("[Register] Registering email: %s", req.Email)

	existing, err := c.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		ctx.JSON(http.StatusConflict, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Email already in use",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("[Register] Failed to hash password: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Failed to process registration",
		})
		return
	}

	user := &models.User{
		Name:          req.Name,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		Age:           req.Age,
		Gender:        req.Gender,
		Weight:        req.Weight,
		Height:        req.Height,
		Goal:          req.Goal,
		ActivityLevel: req.ActivityLevel,
		CreatedAt:     time.Now(),
		Role:          "user",
	}

	if err := c.userRepo.Create(user); err != nil {
		log.Printf("[Register] Failed to create user: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Failed to create user account",
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.LoginResponseDTO{
		Status:  "success",
		Message: "User registered successfully",
		Data: dto.AuthResponseDataDTO{
			User: toAuthUserResponse(user),
		},
	})
}

// Login authenticates a user and sets HttpOnly JWT cookies on success.
// The response body contains user info and the access token expiry time.
// The actual tokens are NOT in the response body — they are in cookies.
func (c *UserAuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[Login] Invalid request: %v", err)
		ctx.JSON(http.StatusBadRequest, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Invalid request format",
		})
		return
	}

	log.Printf("[Login] Attempt for: %s", req.Email)

	user, err := c.userRepo.FindByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	if err := utils.ComparePasswords(user.PasswordHash, req.Password); err != nil {
		log.Printf("[Login] Password mismatch for: %s", req.Email)
		ctx.JSON(http.StatusUnauthorized, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	tokenPair, err := c.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("[Login] Token generation failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Failed to generate authentication token",
		})
		return
	}

	c.setTokenCookies(ctx, tokenPair)

	ctx.JSON(http.StatusOK, dto.LoginResponseDTO{
		Status:  "success",
		Message: "Login successful",
		Data: dto.AuthResponseDataDTO{
			User:      toAuthUserResponse(user),
			ExpiresIn: tokenPair.ExpiresIn,
		},
	})
}

// Logout clears both auth cookies, effectively ending the session.
// No server-side token store to clean up — the JWT simply becomes unused.
func (c *UserAuthController) Logout(ctx *gin.Context) {
	c.clearTokenCookies(ctx)
	ctx.JSON(http.StatusOK, dto.LoginResponseDTO{
		Status:  "success",
		Message: "Logged out successfully",
	})
}

// Refresh reads the refresh_token cookie, validates it, and issues a new token pair.
// The new access_token and refresh_token cookies replace the old ones.
func (c *UserAuthController) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(auth.RefreshTokenCookie)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Refresh token not found",
		})
		return
	}

	tokenPair, err := c.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		log.Printf("[Refresh] Invalid refresh token: %v", err)
		ctx.JSON(http.StatusUnauthorized, dto.LoginResponseDTO{
			Status:  "error",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	c.setTokenCookies(ctx, tokenPair)

	ctx.JSON(http.StatusOK, dto.LoginResponseDTO{
		Status:  "success",
		Message: "Token refreshed successfully",
		Data:    gin.H{"expires_in": tokenPair.ExpiresIn},
	})
}
