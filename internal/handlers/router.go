package handlers

import (
	"net/http"
	"shahbaz-trades-fno/internal/stoxkart"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, stoxClient *stoxkart.Client, whitelistedIP string) {
	// Register IP endpoint
	huma.Register(api, huma.Operation{
		OperationID: "get-ip",
		Method:      http.MethodGet,
		Path:        "/ip",
		Summary:     "Get whitelisted IP",
		Description: "Returns the public IP of the server that needs to be whitelisted in Stoxkart.",
		Tags:        []string{"System"},
	}, HandleGetIP(whitelistedIP))

	// Register Place Order endpoint
	huma.Register(api, huma.Operation{
		OperationID: "place-order",
		Method:      http.MethodPost,
		Path:        "/order/place",
		Summary:     "Place a buy/sell order in FnO",
		Description: "Places an order using Stoxkart Superr API",
		Tags:        []string{"Orders"},
	}, HandlePlaceOrder(stoxClient))
}
