package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// GreetingOutput represents the response body.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!"`
	}
}

func main() {

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err == nil {
		defer resp.Body.Close()
		ip, _ := io.ReadAll(resp.Body)
		fmt.Printf("\n-----------------------------------------\n")
		fmt.Printf("WHITELIST THIS IP: %s\n", string(ip))
		fmt.Printf("-----------------------------------------\n\n")
	}

	// Create a new Gin router
	router := gin.Default()

	// Create a Huma API atop the Gin router
	config := huma.DefaultConfig("Hello World API", "1.0.0")
	api := humagin.New(router, config)

	// Register a GET /hello/{name} endpoint
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/hello/{name}",
		Summary:     "Get a greeting",
		Description: "Returns a hello greeting for the specified name.",
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxlength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	router.Run(":8090")
}
