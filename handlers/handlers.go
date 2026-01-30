package handlers

import (
	"context"
	"net/http"

	"shahbaztradesfno/middleware"
	"shahbaztradesfno/models"
	"shahbaztradesfno/service"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, authService *service.AuthService, tradeService *service.OrderService) {
	authMiddleware := middleware.AuthMiddleware()
	// --- Route 1: Login ---
	huma.Register(api, huma.Operation{
		OperationID: "login-step-1",
		Method:      http.MethodPost,
		Path:        "/api/auth/typea/login",
		Summary:     "Initiate Login",
		Middlewares: huma.Middlewares{authMiddleware},
	}, func(ctx context.Context, input *struct{ Body models.LoginInput }) (*struct{ Body models.APIResponse }, error) {
		resp, err := authService.Login(input.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("Login Request Failed", err)
		}
		if resp.Status == "error" {
			return nil, huma.Error401Unauthorized(resp.Message)
		}
		return &struct{ Body models.APIResponse }{Body: *resp}, nil
	})

	// --- Route 2: Verify OTP ---
	huma.Register(api, huma.Operation{
		OperationID: "verify-otp",
		Method:      http.MethodPost,
		Path:        "/api/auth/typea/verify",
		Summary:     "Verify OTP",
		Middlewares: huma.Middlewares{authMiddleware},
	}, func(ctx context.Context, input *struct{ Body models.VerifyInput }) (*struct{ Body models.APIResponse }, error) {
		resp, err := authService.Verify(input.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("Verification Request Failed", err)
		}
		if resp.Status == "error" {
			if resp.Message == "No pending login session found for this user." {
				return nil, huma.Error404NotFound(resp.Message)
			}
			return nil, huma.Error401Unauthorized(resp.Message)
		}
		return &struct{ Body models.APIResponse }{Body: *resp}, nil
	})

	// --- Route 3: Place Order ---
	huma.Register(api, huma.Operation{
		OperationID: "place-order",
		Method:      http.MethodPost,
		Path:        "/api/trade/order",
		Summary:     "Place Order",
		Middlewares: huma.Middlewares{authMiddleware},
	}, func(ctx context.Context, input *struct{ Body models.OrderInput }) (*struct{ Body models.APIResponse }, error) {
		resp, err := tradeService.PlaceOrder(input.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("Order Request Failed", err)
		}
		if resp.Status == "error" {
			return nil, huma.Error401Unauthorized(resp.Message)
		}
		return &struct{ Body models.APIResponse }{Body: *resp}, nil
	})

	// --- Route 4: Get IP ---
	huma.Register(api, huma.Operation{
		OperationID: "get-ip",
		Method:      http.MethodGet,
		Path:        "/api/ip",
		Summary:     "Get Public IP",
	}, func(ctx context.Context, input *struct{}) (*struct{ Body models.APIResponse }, error) {
		ip, err := authService.GetIP()
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to get IP", err)
		}
		return &struct{ Body models.APIResponse }{Body: models.APIResponse{
			Status:  "success",
			Message: "IP fetched successfully",
			Data:    map[string]string{"ip": ip},
		}}, nil
	})
}
