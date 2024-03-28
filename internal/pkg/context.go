package pkg

import "context"

const (
	ApiKeyContextKey = "API_KEY"
)

func GetApiKeyFromContext(ctx context.Context) string {
	u, ok := ctx.Value(ApiKeyContextKey).(string)
	if !ok {
		return ""
	}
	return u
}
