package service

import (
	"fmt"

	"shahbaztradesfno/cache"
	"shahbaztradesfno/models"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) PlaceOrder(input models.OrderInput) (*models.APIResponse, error) {
	val, ok := cache.SessionCache.Get(input.Username)
	if !ok {
		return nil, fmt.Errorf("session not found for user %s", input.Username)
	}

	userSvc := val.(*UserService)
	if userSvc == nil {
		return nil, fmt.Errorf("user service not found for user %s", input.Username)
	}

	variety := input.Variety
	if variety == "" {
		variety = "regular"
	}

	resp, err := userSvc.client.R().
		SetHeader("X-Mirae-Version", "1").
		SetHeader("Authorization", fmt.Sprintf("token %s:%s", userSvc.ApiKey, userSvc.AccessToken)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"tradingsymbol":    input.Symbol,
			"exchange":         input.Exchange,
			"transaction_type": input.Side,
			"order_type":       input.Type,
			"quantity":         input.Qty,
			"product":          input.Product,
			"validity":         input.Validity,
			"price":            input.Price,
		}).
		Post("https://api.mstock.trade/openapi/typea/orders/" + variety)

	if err != nil {
		log.Error().Err(err).Str("username", input.Username).Str("symbol", input.Symbol).Msg("MStock order placement connection failed")
		return nil, err
	}

	var result map[string]any
	if err := sonic.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Str("username", input.Username).Msg("Failed to unmarshal MStock order response")
		return nil, err
	}
	log.Info().Str("username", input.Username).Str("symbol", input.Symbol).Str("side", input.Side).Msg("Order placement response received")

	// Check for MStock specific error status
	if status, ok := result["status"].(string); ok && status == "error" {
		msg, _ := result["message"].(string)
		if msg == "" {
			msg = "Order placement failed."
		}
		// Include error_type in the data if available
		log.Warn().Str("username", input.Username).Str("symbol", input.Symbol).Interface("error", result).Msg("Order placement failed by MStock")
		return &models.APIResponse{
			Status:  "error",
			Message: msg,
			Data:    result,
		}, nil
	}

	// For success, return the data part directly
	orderData := result["data"]
	log.Info().Str("username", input.Username).Str("symbol", input.Symbol).Interface("orderData", orderData).Msg("Order placed successfully")
	return &models.APIResponse{
		Status:  "success",
		Message: "Order placed successfully",
		Data:    orderData,
	}, nil
}
