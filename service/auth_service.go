package service

import (
	"shahbaztradesfno/cache"
	"shahbaztradesfno/models"
	"shahbaztradesfno/utils"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
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
		log.Info().Str("ip", s.publicIP).Msg("AuthService initialized with public IP")
	} else {
		log.Error().Err(err).Msg("Failed to fetch public IP during AuthService initialization")
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
		log.Error().Err(err).Str("username", input.Username).Msg("MStock login connection failed")
		return nil, err
	}

	var result map[string]any
	sonic.Unmarshal(resp.Body(), &result)

	if status, ok := result["status"].(string); ok && status == "success" {
		s.mu.Lock()
		cache.PendingLoginCache.Set(input.Username, input.APIKey, 5*time.Minute)
		s.mu.Unlock()

		log.Info().Str("username", input.Username).Msg("MStock login initiated (Step 1 success)")
		return &models.APIResponse{
			Status:  "success",
			Message: "OTP sent. Please verify with OTP and Username.",
		}, nil
	}

	msg, _ := result["message"].(string)
	if msg == "" {
		msg = "m.Stock login failed without a specific reason."
	}
	log.Warn().Str("username", input.Username).Str("reason", msg).Msg("MStock login failed")

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
		log.Error().Err(err).Str("username", input.Username).Msg("MStock verification connection failed")
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

		log.Info().Str("username", input.Username).Msg("MStock session established successfully")
		return &models.APIResponse{
			Status:  "success",
			Message: "Session established successfully",
		}, nil
	}

	msg, _ := result["message"].(string)
	if msg == "" {
		msg = "OTP verification failed."
	}
	log.Warn().Str("username", input.Username).Str("reason", msg).Msg("MStock verification failed")
	return &models.APIResponse{Status: "error", Message: msg}, nil
}

func (s *AuthService) GetIP() (string, error) {
	return s.publicIP, nil
}
