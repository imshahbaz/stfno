package middleware

import (
	"os"

	"github.com/danielgtaylor/huma/v2"
)

var apiSecret string

func init() {
	apiSecret = os.Getenv("API_SECRET")
	if apiSecret == "" {
		apiSecret = "shahbaz-trades-secret"
	}
}

func AuthMiddleware() func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		apiKey := ctx.Header("X-Auth-Key")

		if apiKey != apiSecret {
			ctx.SetStatus(401)
			ctx.BodyWriter().Write([]byte(`{"status":"error","message":"Unauthorized: Invalid or missing API key"}`))
			return
		}

		next(ctx)
	}
}
