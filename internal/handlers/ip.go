package handlers

import (
	"context"
)

type GetIPInput struct{}

type GetIPOutput struct {
	Body struct {
		IP string `json:"ip" doc:"The public IP of the server"`
	}
}

func HandleGetIP(whitelistedIP string) func(ctx context.Context, input *GetIPInput) (*GetIPOutput, error) {
	return func(ctx context.Context, input *GetIPInput) (*GetIPOutput, error) {
		resp := &GetIPOutput{}
		resp.Body.IP = whitelistedIP
		return resp, nil
	}
}
