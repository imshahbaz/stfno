package main

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"shahbaztradesfno/handlers"
	"shahbaztradesfno/service"
)

func main() {
	authService := service.NewAuthService()
	orderService := service.NewOrderService()

	router := gin.Default()
	config := huma.DefaultConfig("Shahbaz Trades FnO API", "1.0.0")
	api := humagin.New(router, config)

	handlers.RegisterRoutes(api, authService, orderService)

	router.Run(":3000")
}
