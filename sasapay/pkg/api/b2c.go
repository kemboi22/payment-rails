package api

import (
	"encoding/json"
	"fmt"
)

type C2BRequest struct {
	MerchantCode     string `json:"MerchantCode"`
	NetworkCode      string `json:"NetworkCode"`
	TransactionFee   int    `json:"Transaction Fee"`
	Currency         string `json:"Currency"`
	Amount           string `json:"Amount"`
	CallBackURL      string `json:"CallBackURL"`
	PhoneNumber      string `json:"PhoneNumber"`
	TransactionDesc  string `json:"TransactionDesc"`
	AccountReference string `json:"AccountReference"`
}
type C2BResponse struct {
	Status               bool   `json:"status"`
	Detail               string `json:"detail"`
	PaymentGateway       string `json:"PaymentGateway"`
	MerchantRequestID    string `json:"MerchantRequestID"`
	CheckoutRequestID    string `json:"CheckoutRequestID"`
	TransactionReference string `json:"TransactionReference"`
	ResponseCode         string `json:"ResponseCode"`
	ResponseDescription  string `json:"ResponseDescription"`
	CustomerMessage      string `json:"CustomerMessage"`
}

type ProcessPaymentRequest struct {
	CheckoutRequestID string `json:"CheckoutRequestID"`
	MerchantCode      string `json:"MerchantCode"`
	VerificationCode  string `json:"VerificationCode"`
}
type ProcessPaymentReponse struct {
	Status bool   `json:"status"`
	Detail string `json:"detail"`
}

const (
	RequestPaymentURL = "/v1/payments/request-payment/"
	ProcessPaymentURL = "/v1/payments/process-payment/"
)

// Request Payment From Sasa Pay User
// SasaPay User
// NetworkCode = 0
// Other Users NetworkCode = 63902(M-Pesa) 63903 (Airtel Money) 63907(T-Kash)
func (client *Client) RequestPayment(req *C2BRequest) (*C2BResponse, error) {
	resp, err := client.MakeRequest("post", RequestPaymentURL, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to request payment: %w", err)
	}
	var response *C2BResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the json: %w", err)
	}
	return response, nil
}

func (client *Client) ProcessPayment(req *ProcessPaymentRequest) (*ProcessPaymentReponse, error) {
	resp, err := client.MakeRequest("post", ProcessPaymentURL, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}
	var response *ProcessPaymentReponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return response, nil
}
