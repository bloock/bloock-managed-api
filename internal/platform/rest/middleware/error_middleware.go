package middleware

import (
	"errors"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return errorMiddleware(gin.ErrorTypeAny)
}

func errorMiddleware(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError api_error.APIError

			switch v := err.(type) {
			case *api_error.APIError:
				parsedError = *v
			default:
				parsedError = *api_error.NewInternalServerAPIError(errors.New("Internal Server Error"))
			}

			c.AbortWithStatusJSON(parsedError.Status, parsedError)
			return
		}

	}
}
