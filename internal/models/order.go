package models

type OrderRequest struct {
	Exchange          string  `json:"exchange" doc:"NFO for FnO"`
	SymbolToken       string  `json:"symboltoken" doc:"The specific token for the Option/Future"`
	TransactionType   string  `json:"transactiontype" doc:"BUY or SELL"`
	Quantity          int     `json:"quantity" doc:"Lot size (e.g., 75 for Nifty)"`
	DisclosedQuantity int     `json:"disclosedquantity" doc:"Usually 0"`
	Price             float64 `json:"price" doc:"Price for LIMIT orders"`
	Type              string  `json:"type" doc:"LIMIT or MARKET"`
	ProductType       string  `json:"producttype" doc:"NRML (Carry forward) or MIS (Intraday)"`
	Validity          string  `json:"validity" doc:"DAY"`
	OrderTag          string  `json:"ordertag" doc:"Optional tag for tracking"`
}

type OrderResponse struct {
	Status string `json:"status"`
	Data   struct {
		Message string `json:"message"`
		OrderID string `json:"order_id"`
	} `json:"data"`
}
