package main

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"shahbaztradesfno/handlers"
	"shahbaztradesfno/service"

	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	port string
	env  string
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	port = os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	env = os.Getenv("ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func main() {
	authService := service.NewAuthService()
	orderService := service.NewOrderService()

	router := gin.Default()
	config := huma.DefaultConfig("Shahbaz Trades FnO API", "1.0.0")
	api := humagin.New(router, config)

	handlers.RegisterRoutes(api, authService, orderService)

	log.Info().Msg("Starting Shahbaz Trades FnO API on " + port)
	router.Run(":" + port)
}
