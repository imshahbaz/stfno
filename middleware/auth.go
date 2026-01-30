package middleware

import (
	"os"

	"github.com/danielgtaylor/huma/v2"
)

var apiSecret string
var authEnabled bool

func init() {
	apiSecret = os.Getenv("API_SECRET")
	if apiSecret == "" {
		apiSecret = "shahbaz-trades-secret"
	}
	authEnabledEnv := os.Getenv("AUTH_ENABLED")
	if authEnabledEnv == "" {
		authEnabled = true
	} else {
		authEnabled = authEnabledEnv == "true"
	}
}

func AuthMiddleware() func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		apiKey := ctx.Header("X-Auth-Key")

		if authEnabled && apiKey != apiSecret {
			ctx.SetStatus(401)
			ctx.BodyWriter().Write([]byte(`{"status":"error","message":"Unauthorized: Invalid or missing API key"}`))
			return
		}

		next(ctx)
	}
}
