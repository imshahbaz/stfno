package service

import (
	"shahbaztradesfno/cache"
	"shahbaztradesfno/models"
	"shahbaztradesfno/utils"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type AuthService struct {
	mu       sync.RWMutex
	client   *resty.Client
	publicIP string
}

func NewAuthService() *AuthService {
	client := resty.New()
	s := &AuthService{
		client: client,
	}
	// Fetch IP once at startup
	resp, err := client.R().Get("https://api.ipify.org")
	if err == nil {
		s.publicIP = string(resp.Body())
	}
	return s
}

func (s *AuthService) Login(input models.LoginInput) (*models.APIResponse, error) {
	resp, err := s.client.R().
		SetHeader("X-Mirae-Version", "1").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"username": input.Username,
			"password": input.Password,
		}).
		Post("https://api.mstock.trade/openapi/typea/connect/login")

	if err != nil {
		return nil, err
	}

	var result map[string]any
	sonic.Unmarshal(resp.Body(), &result)

	if status, ok := result["status"].(string); ok && status == "success" {
		s.mu.Lock()
		cache.PendingLoginCache.Set(input.Username, input.APIKey, 5*time.Minute)
		s.mu.Unlock()

		return &models.APIResponse{
			Status:  "success",
			Message: "OTP sent. Please verify with OTP and Username.",
		}, nil
	}

	msg, _ := result["message"].(string)
	if msg == "" {
		msg = "m.Stock login failed without a specific reason."
	}

	return &models.APIResponse{Status: "error", Message: msg}, nil
}

func (s *AuthService) Verify(input models.VerifyInput) (*models.APIResponse, error) {
	pending, exists := cache.PendingLoginCache.Get(input.Username)

	if !exists {
		return &models.APIResponse{Status: "error", Message: "No pending login session found for this user."}, nil
	}

	resp, err := s.client.R().
		SetHeader("X-Mirae-Version", "1").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"api_key":       pending.(string),
			"request_token": input.OTP,
		}).
		Post("https://api.mstock.trade/openapi/typea/session/token")

	if err != nil {
		return nil, err
	}

	var result map[string]any
	sonic.Unmarshal(resp.Body(), &result)

	if status, ok := result["status"].(string); ok && status == "success" {
		data, _ := result["data"].(map[string]any)

		s.mu.Lock()
		userSvc := NewUserService(pending.(string), data["access_token"].(string))
		cache.SessionCache.Set(input.Username, userSvc, utils.GetDurationToMidnightIST())
		s.mu.Unlock()

		return &models.APIResponse{
			Status:  "success",
			Message: "Session established successfully",
		}, nil
	}

	msg, _ := result["message"].(string)
	if msg == "" {
		msg = "OTP verification failed."
	}
	return &models.APIResponse{Status: "error", Message: msg}, nil
}

func (s *AuthService) GetIP() (string, error) {
	return s.publicIP, nil
}
