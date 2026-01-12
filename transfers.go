package wise

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// TransfersService handles transfer-related API calls.
type TransfersService struct {
	client *Client
}

// Transfer represents a money transfer.
type Transfer struct {
	ID                    int64           `json:"id"`
	User                  int64           `json:"user"`
	TargetAccount         int64           `json:"targetAccount"`
	SourceAccount         int64           `json:"sourceAccount,omitempty"`
	Quote                 string          `json:"quote,omitempty"`
	QuoteUUID             string          `json:"quoteUuid,omitempty"`
	Status                TransferStatus  `json:"status"`
	Rate                  float64         `json:"rate"`
	Reference             string          `json:"reference,omitempty"`
	Created               Timestamp       `json:"created"`
	Business              int64           `json:"business,omitempty"`
	TransferRequest       int64           `json:"transferRequest,omitempty"`
	Details               TransferDetails `json:"details"`
	HasActiveIssues       bool            `json:"hasActiveIssues"`
	SourceCurrency        Currency        `json:"sourceCurrency"`
	SourceValue           float64         `json:"sourceValue"`
	TargetCurrency        Currency        `json:"targetCurrency"`
	TargetValue           float64         `json:"targetValue"`
	CustomerTransactionID string          `json:"customerTransactionId,omitempty"`
}

// TransferDetails represents additional details of a transfer.
type TransferDetails struct {
	Reference       string `json:"reference,omitempty"`
	TransferPurpose string `json:"transferPurpose,omitempty"`
	SourceOfFunds   string `json:"sourceOfFunds,omitempty"`
}

// CreateTransferRequest represents the request to create a transfer.
type CreateTransferRequest struct {
	TargetAccount         int64           `json:"targetAccount"`
	QuoteUUID             string          `json:"quoteUuid"`
	CustomerTransactionID string          `json:"customerTransactionId"`
	Details               TransferDetails `json:"details,omitempty"`
}

// FundTransferRequest represents the request to fund a transfer.
type FundTransferRequest struct {
	Type string `json:"type"` // BALANCE
}

// TransferIssue represents an issue with a transfer.
type TransferIssue struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// ListTransfersParams represents the parameters for listing transfers.
type ListTransfersParams struct {
	ProfileID int64
	Status    TransferStatus
	Limit     int
	Offset    int
	CreatedDateStart string // ISO 8601 format
	CreatedDateEnd   string // ISO 8601 format
}

// Create creates a new transfer.
// POST /v1/transfers
func (s *TransfersService) Create(ctx context.Context, req *CreateTransferRequest) (*Transfer, error) {
	var transfer Transfer
	err := s.client.Post(ctx, "/v1/transfers", req, &transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// Get retrieves a transfer by ID.
// GET /v1/transfers/{transferId}
func (s *TransfersService) Get(ctx context.Context, transferID int64) (*Transfer, error) {
	var transfer Transfer
	path := fmt.Sprintf("/v1/transfers/%d", transferID)
	err := s.client.Get(ctx, path, nil, &transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// List returns transfers based on filters.
// GET /v1/transfers
func (s *TransfersService) List(ctx context.Context, params *ListTransfersParams) ([]Transfer, error) {
	query := url.Values{}
	if params != nil {
		if params.ProfileID > 0 {
			query.Set("profile", strconv.FormatInt(params.ProfileID, 10))
		}
		if params.Status != "" {
			query.Set("status", string(params.Status))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.CreatedDateStart != "" {
			query.Set("createdDateStart", params.CreatedDateStart)
		}
		if params.CreatedDateEnd != "" {
			query.Set("createdDateEnd", params.CreatedDateEnd)
		}
	}

	var transfers []Transfer
	err := s.client.Get(ctx, "/v1/transfers", query, &transfers)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

// Cancel cancels a transfer.
// PUT /v1/transfers/{transferId}/cancel
func (s *TransfersService) Cancel(ctx context.Context, transferID int64) (*Transfer, error) {
	var transfer Transfer
	path := fmt.Sprintf("/v1/transfers/%d/cancel", transferID)
	err := s.client.Put(ctx, path, nil, &transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// Fund funds a transfer from a balance.
// POST /v3/profiles/{profileId}/transfers/{transferId}/payments
func (s *TransfersService) Fund(ctx context.Context, profileID, transferID int64) (*Transfer, error) {
	req := FundTransferRequest{Type: "BALANCE"}
	var transfer Transfer
	path := fmt.Sprintf("/v3/profiles/%d/transfers/%d/payments", profileID, transferID)
	err := s.client.Post(ctx, path, req, &transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetIssues retrieves issues for a transfer.
// GET /v1/transfers/{transferId}/issues
func (s *TransfersService) GetIssues(ctx context.Context, transferID int64) ([]TransferIssue, error) {
	var issues []TransferIssue
	path := fmt.Sprintf("/v1/transfers/%d/issues", transferID)
	err := s.client.Get(ctx, path, nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

// GetDeliveryTime gets the estimated delivery time for a transfer.
// GET /v1/delivery-estimates/{transferId}
func (s *TransfersService) GetDeliveryTime(ctx context.Context, transferID int64) (*Timestamp, error) {
	var result struct {
		EstimatedDeliveryDate Timestamp `json:"estimatedDeliveryDate"`
	}
	path := fmt.Sprintf("/v1/delivery-estimates/%d", transferID)
	err := s.client.Get(ctx, path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.EstimatedDeliveryDate, nil
}
