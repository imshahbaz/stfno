package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"shahbaz-trades-fno/internal/handlers"
	"shahbaz-trades-fno/internal/stoxkart"

	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Fetch and store Whitelist IP
	var whitelistedIP string
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err == nil {
		defer resp.Body.Close()
		ip, _ := io.ReadAll(resp.Body)
		whitelistedIP = string(ip)
		fmt.Printf("\n-----------------------------------------\n")
		fmt.Printf("WHITELIST THIS IP: %s\n", whitelistedIP)
		fmt.Printf("-----------------------------------------\n\n")
	}

	// Create a new Gin router
	router := gin.Default()

	// Create a Huma API atop the Gin router
	config := huma.DefaultConfig("Stoxkart FnO Trading API", "1.0.0")
	config.Formats["application/json"] = huma.Format{
		Marshal: func(w io.Writer, v any) error {
			return sonic.ConfigDefault.NewEncoder(w).Encode(v)
		},
		Unmarshal: sonic.Unmarshal,
	}
	api := humagin.New(router, config)

	// Initialize Stoxkart Client
	stoxClient := stoxkart.NewClient()

	// Register all routes
	handlers.RegisterRoutes(api, stoxClient, whitelistedIP)

	// Start the server
	router.Run(":8090")
}
