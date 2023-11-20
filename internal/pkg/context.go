package pkg

import "context"

const (
	ApiKeyContextKey = "API_KEY"
	EnvContextKey    = "ENV"
)

func GetApiKeyFromContext(ctx context.Context) string {
	u, ok := ctx.Value(ApiKeyContextKey).(string)
	if !ok {
		return ""
	}
	return u
}

func GetEnvFromContext(ctx context.Context) *string {
	u, ok := ctx.Value(EnvContextKey).(string)
	if !ok {
		return nil
	}
	return &u
}
