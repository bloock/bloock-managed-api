package middleware

import (
	"strings"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/pkg"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Configuration.Auth.Secret != "" {
			secret := getBearerAuth(c)
			if secret != config.Configuration.Auth.Secret {
				_ = c.Error(api_error.NewUnauthorizedAPIError("invalid secret provided"))
				c.Abort()
				return
			}
		}

		apiKey := getApiKey(c)
		if apiKey != "" {
			c.Set(pkg.ApiKeyContextKey, apiKey)
		} else {
			c.Set(pkg.ApiKeyContextKey, config.Configuration.Bloock.ApiKey)
		}

		environment := getEnvironment(c)
		if environment != "" {
			c.Set(pkg.EnvContextKey, environment)
		}

		c.Next()
	}
}

func getBearerAuth(c *gin.Context) string {
	authorizationHeader := c.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return authorizationHeader
	}

	splitToken := strings.Split(authorizationHeader, "Bearer")
	if len(splitToken) != 2 {
		return ""
	}

	return strings.TrimSpace(splitToken[1])
}

func getApiKey(c *gin.Context) string {
	return c.Request.Header.Get("X-Api-Key")
}

func getEnvironment(c *gin.Context) string {
	return c.Request.Header.Get("Environment")
}
