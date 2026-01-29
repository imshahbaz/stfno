package stoxkart

import (
	"fmt"
	"shahbaz-trades-fno/internal/models"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	restyClient *resty.Client
}

func NewClient() *Client {
	c := resty.New()

	// Use sonic for JSON marshaling/unmarshaling
	c.JSONMarshal = sonic.Marshal
	c.JSONUnmarshal = sonic.Unmarshal

	return &Client{
		restyClient: c,
	}
}

func (c *Client) PlaceOrder(jwtToken string, order models.OrderRequest) (*models.OrderResponse, error) {
	const orderUrl = "https://swapi.stoxkart.com/v1/order/placeOrder"

	var orderResponse models.OrderResponse
	resp, err := c.restyClient.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader("Content-Type", "application/json").
		SetBody(order).
		SetResult(&orderResponse).
		Post(orderUrl)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}

	return &orderResponse, nil
}
