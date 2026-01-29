package handlers

import (
	"context"
	"shahbaz-trades-fno/internal/models"
	"shahbaz-trades-fno/internal/stoxkart"
)

type PlaceOrderInput struct {
	Authorization string `header:"Authorization" doc:"Bearer [your_jwt_token]"`
	Body          models.OrderRequest
}

type PlaceOrderOutput struct {
	Body models.OrderResponse
}

func HandlePlaceOrder(client *stoxkart.Client) func(ctx context.Context, input *PlaceOrderInput) (*PlaceOrderOutput, error) {
	return func(ctx context.Context, input *PlaceOrderInput) (*PlaceOrderOutput, error) {
		// Extract token if "Bearer " prefix is present
		token := input.Authorization
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		resp, err := client.PlaceOrder(token, input.Body)
		if err != nil {
			return nil, err
		}

		return &PlaceOrderOutput{
			Body: *resp,
		}, nil
	}
}
