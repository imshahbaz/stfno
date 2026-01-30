package service

import "github.com/go-resty/resty/v2"

type UserService struct {
	client      *resty.Client
	ApiKey      string
	AccessToken string
}

func NewUserService(apiKey string, accessToken string) *UserService {
	return &UserService{
		client:      resty.New(),
		ApiKey:      apiKey,
		AccessToken: accessToken,
	}
}
