package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First, try to get token from HttpOnly cookie (preferred for security)
		if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// Fall back to Authorization header for backwards compatibility and API clients
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				slog.Warn("Auth failed: no token or authorization header", "ip", c.ClientIP(), "path", c.Request.URL.Path)
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Code:    apperror.CodeUnauthorized,
					Message: "authorization required",
				})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				slog.Warn("Auth failed: invalid authorization header format", "ip", c.ClientIP(), "path", c.Request.URL.Path)
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Code:    apperror.CodeUnauthorized,
					Message: "invalid authorization header format",
				})
				c.Abort()
				return
			}
			tokenString = parts[1]
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			slog.Warn("Auth failed: invalid token", "ip", c.ClientIP(), "path", c.Request.URL.Path, "error", err)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    apperror.CodeUnauthorized,
				Message: "invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Warn("Auth failed: invalid token claims", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    apperror.CodeUnauthorized,
				Message: "invalid token claims",
			})
			c.Abort()
			return
		}

		// Verify this is an access token (not a refresh token)
		// For backwards compatibility, accept tokens without type claim
		if tokenType, exists := claims["type"]; exists {
			if tokenType != "access" {
				slog.Warn("Auth failed: invalid token type", "ip", c.ClientIP(), "path", c.Request.URL.Path, "type", tokenType)
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Code:    apperror.CodeUnauthorized,
					Message: "invalid token type",
				})
				c.Abort()
				return
			}
		}

		// JWT numbers are parsed as float64, convert to uint
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			slog.Warn("Auth failed: invalid user id in token", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    apperror.CodeUnauthorized,
				Message: "invalid user id in token",
			})
			c.Abort()
			return
		}

		c.Set("userID", uint(userIDFloat))
		c.Set("userEmail", claims["email"])
		c.Next()
	}
}
