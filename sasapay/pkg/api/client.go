package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

type Environment string

const (
	SANDBOX    Environment = "sandbox"
	PRODUCTION Environment = "production"
	AuthURL                = "/v1/auth/token/?grant_type=client_credentials"
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

func (client *Client) GetAuthToken() (string, error) {
	if token, found := client.Cache.Get("sasapay_auth_token"); found {
		return token.(string), nil
	}
	url := client.BaseURL + AuthURL
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	auth := base64.StdEncoding.EncodeToString([]byte(client.ClientID + ":" + client.ClientSecret))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to execute auth request with status : %s", resp.Status)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	expiresIn := 3600
	client.Cache.Set("sasapay_auth_token", authResponse.AccessToken, time.Duration(expiresIn)*time.Second)

	return authResponse.AccessToken, nil
}

func (client *Client) MakeRequest(method, url string, payload interface{}) ([]byte, error) {
	token, err := client.GetAuthToken()
	if err != nil {
		return nil, err
	}
	var reqBody []byte
	if payload != nil {
		reqBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}
	req, err := http.NewRequest(method, client.BaseURL+url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization ", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %s", resp.Status)
	}

	respBody, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	var respBody bytes.Buffer
	_, err := respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return respBody.Bytes(), nil
}
