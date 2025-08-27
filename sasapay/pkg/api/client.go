package api

import (
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

type Environment string

const (
	SANDBOX    Environment = "sandbox"
	PRODUCTION Environment = "production"
	AuthURL                = "auth/token/?grant_type=client_credentials"
)

type AuthResponse struct {
	Status      string `json:"status"`
	Detail      string `json:"detail"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	Environment  Environment
}

type Client struct {
	ClientID     string
	ClientSecret string
	Environment  Environment
	HTTPClient   *http.Client
	Cache        *cache.Cache
	BaseURL      string
}

func NewClient(credentials ClientCredentials) *Client {
	c := cache.New(50*time.Minute, 10*time.Minute)
	var baseURL string
	if credentials.Environment == "production" {
		baseURL = "https://api.sasapay.app/api"
	} else {
		baseURL = "https://sandbox.sasapay.app/api"
	}
	return &Client{
		ClientID:     credentials.ClientID,
		ClientSecret: credentials.ClientSecret,
		Cache:        c,
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		BaseURL:      baseURL,
	}
}
