package middleware

import (
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

var apiSecret string
var authEnabled bool

func init() {
	apiSecret = os.Getenv("API_SECRET")
	if apiSecret == "" {
		apiSecret = "shahbaz-trades-secret"
	}

	authEnabledEnv := os.Getenv("AUTH_ENABLED")
	authEnabled = authEnabledEnv == "true"
	log.Info().Bool("authEnabled", authEnabled).Msg("Auth enabled")

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
