package models

type LoginInput struct {
	Username string `json:"username" form:"username" doc:"m.Stock Client Code"`
	Password string `json:"password" form:"password" doc:"m.Stock Password"`
	APIKey   string `json:"apiKey" form:"apiKey" doc:"m.Stock Type A API Key"`
}

type VerifyInput struct {
	Username string `json:"username" form:"username" doc:"m.Stock Client Code"`
	OTP      string `json:"otp" form:"otp" doc:"6-digit SMS OTP"`
}

type OrderInput struct {
	Username string `json:"username" form:"username"`
	Symbol   string `json:"symbol" form:"symbol" doc:"Trading Symbol e.g. INFY"`
	Exchange string `json:"exchange" form:"exchange" doc:"NSE or BSE"`
	Side     string `json:"side" form:"side" doc:"BUY or SELL"`
	Type     string `json:"type" form:"type" doc:"LIMIT or MARKET"`
	Qty      string `json:"qty" form:"qty"`
	Product  string `json:"product" form:"product" doc:"MIS, DELIVERY, etc."`
	Validity string `json:"validity" form:"validity" doc:"DAY, etc."`
	Price    string `json:"price" form:"price"`
	Variety  string `json:"variety" form:"variety" default:"regular"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Session struct {
	AccessToken string
	APIKey      string
}
