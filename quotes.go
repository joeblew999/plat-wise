package wise

import (
	"context"
	"fmt"
)

// QuotesService handles quote-related API calls.
type QuotesService struct {
	client *Client
}

// Quote represents a quote for a money transfer.
type Quote struct {
	ID                   string        `json:"id"`
	SourceCurrency       Currency      `json:"sourceCurrency"`
	TargetCurrency       Currency      `json:"targetCurrency"`
	SourceAmount         float64       `json:"sourceAmount,omitempty"`
	TargetAmount         float64       `json:"targetAmount,omitempty"`
	PayOut               string        `json:"payOut,omitempty"`
	Rate                 float64       `json:"rate"`
	CreatedTime          Timestamp     `json:"createdTime"`
	User                 int64         `json:"user"`
	Profile              int64         `json:"profile"`
	RateType             string        `json:"rateType,omitempty"`
	RateExpirationTime   Timestamp     `json:"rateExpirationTime"`
	GuaranteedTargetAmount bool        `json:"guaranteedTargetAmount,omitempty"`
	ProvidedAmountType   string        `json:"providedAmountType,omitempty"`
	PaymentOptions       []PaymentOption `json:"paymentOptions,omitempty"`
	Status               string        `json:"status,omitempty"`
	ExpirationTime       Timestamp     `json:"expirationTime,omitempty"`
}

// PaymentOption represents a payment option for a quote.
type PaymentOption struct {
	FormattedEstimatedDelivery string     `json:"formattedEstimatedDelivery,omitempty"`
	EstimatedDelivery          Timestamp  `json:"estimatedDelivery,omitempty"`
	AllowedProfileTypes        []string   `json:"allowedProfileTypes,omitempty"`
	PayInProduct               string     `json:"payInProduct,omitempty"`
	FeePercentage              float64    `json:"feePercentage,omitempty"`
	EstimatedDeliveryDelays    []string   `json:"estimatedDeliveryDelays,omitempty"`
	Fee                        Money      `json:"fee,omitempty"`
	SourceAmount               float64    `json:"sourceAmount,omitempty"`
	TargetAmount               float64    `json:"targetAmount,omitempty"`
	PayIn                      string     `json:"payIn,omitempty"`
	Disabled                   bool       `json:"disabled,omitempty"`
}

// CreateQuoteRequest represents the request to create a quote.
type CreateQuoteRequest struct {
	SourceCurrency     Currency `json:"sourceCurrency"`
	TargetCurrency     Currency `json:"targetCurrency"`
	SourceAmount       *float64 `json:"sourceAmount,omitempty"`
	TargetAmount       *float64 `json:"targetAmount,omitempty"`
	Profile            int64    `json:"profile,omitempty"`
	PayOut             string   `json:"payOut,omitempty"`             // BANK_TRANSFER, BALANCE, etc.
	PreferredPayIn     string   `json:"preferredPayIn,omitempty"`     // BANK_TRANSFER, BALANCE, etc.
}

// UpdateQuoteRequest represents the request to update a quote.
type UpdateQuoteRequest struct {
	SourceAmount   *float64 `json:"sourceAmount,omitempty"`
	TargetAmount   *float64 `json:"targetAmount,omitempty"`
	PayOut         string   `json:"payOut,omitempty"`
	PreferredPayIn string   `json:"preferredPayIn,omitempty"`
}

// Create creates a new quote.
// POST /v3/profiles/{profileId}/quotes
func (s *QuotesService) Create(ctx context.Context, profileID int64, req *CreateQuoteRequest) (*Quote, error) {
	var quote Quote
	path := fmt.Sprintf("/v3/profiles/%d/quotes", profileID)
	err := s.client.Post(ctx, path, req, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

// CreateV2 creates a new quote using the v2 API (simpler, doesn't require profile ID in path).
// POST /v2/quotes
func (s *QuotesService) CreateV2(ctx context.Context, req *CreateQuoteRequest) (*Quote, error) {
	var quote Quote
	err := s.client.Post(ctx, "/v2/quotes", req, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

// Get retrieves a quote by ID.
// GET /v3/profiles/{profileId}/quotes/{quoteId}
func (s *QuotesService) Get(ctx context.Context, profileID int64, quoteID string) (*Quote, error) {
	var quote Quote
	path := fmt.Sprintf("/v3/profiles/%d/quotes/%s", profileID, quoteID)
	err := s.client.Get(ctx, path, nil, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

// GetV2 retrieves a quote by ID using the v2 API.
// GET /v2/quotes/{quoteId}
func (s *QuotesService) GetV2(ctx context.Context, quoteID string) (*Quote, error) {
	var quote Quote
	path := fmt.Sprintf("/v2/quotes/%s", quoteID)
	err := s.client.Get(ctx, path, nil, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

// Update updates an existing quote.
// PATCH /v3/profiles/{profileId}/quotes/{quoteId}
func (s *QuotesService) Update(ctx context.Context, profileID int64, quoteID string, req *UpdateQuoteRequest) (*Quote, error) {
	var quote Quote
	path := fmt.Sprintf("/v3/profiles/%d/quotes/%s", profileID, quoteID)
	err := s.client.Request(ctx, "PATCH", path, nil, req, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}
