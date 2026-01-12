package wise

import (
	"context"
	"fmt"
	"net/url"
)

// BalancesService handles balance-related API calls.
type BalancesService struct {
	client *Client
}

// Balance represents a multi-currency balance.
type Balance struct {
	ID              int64    `json:"id"`
	ProfileID       int64    `json:"profileId,omitempty"`
	Currency        Currency `json:"currency"`
	Amount          Money    `json:"amount"`
	ReservedAmount  Money    `json:"reservedAmount,omitempty"`
	CashAmount      Money    `json:"cashAmount,omitempty"`
	TotalWorth      Money    `json:"totalWorth,omitempty"`
	Type            string   `json:"type,omitempty"` // STANDARD, SAVINGS
	Name            string   `json:"name,omitempty"`
	Icon            string   `json:"icon,omitempty"`
	CreationTime    Timestamp `json:"creationTime,omitempty"`
	ModificationTime Timestamp `json:"modificationTime,omitempty"`
	Visible         bool     `json:"visible"`
}

// BalanceStatement represents a statement entry.
type BalanceStatement struct {
	Type            string    `json:"type"`
	Date            Timestamp `json:"date"`
	Amount          Money     `json:"amount"`
	TotalFees       Money     `json:"totalFees,omitempty"`
	Details         StatementDetails `json:"details,omitempty"`
	ExchangeDetails *ExchangeDetails `json:"exchangeDetails,omitempty"`
	RunningBalance  Money     `json:"runningBalance,omitempty"`
	ReferenceNumber string    `json:"referenceNumber,omitempty"`
}

// StatementDetails contains additional details for a statement entry.
type StatementDetails struct {
	Type            string `json:"type,omitempty"`
	Description     string `json:"description,omitempty"`
	SenderName      string `json:"senderName,omitempty"`
	SenderAccount   string `json:"senderAccount,omitempty"`
	PaymentReference string `json:"paymentReference,omitempty"`
}

// ExchangeDetails contains exchange information for a statement entry.
type ExchangeDetails struct {
	FromAmount   Money   `json:"fromAmount,omitempty"`
	ToAmount     Money   `json:"toAmount,omitempty"`
	Rate         float64 `json:"rate,omitempty"`
}

// ConvertBalanceRequest represents a request to convert between balances.
type ConvertBalanceRequest struct {
	QuoteID string `json:"quoteId"`
}

// ListBalancesParams represents parameters for listing balances.
type ListBalancesParams struct {
	Types []string // STANDARD, SAVINGS
}

// List retrieves all balances for a profile.
// GET /v4/profiles/{profileId}/balances
func (s *BalancesService) List(ctx context.Context, profileID int64, params *ListBalancesParams) ([]Balance, error) {
	query := url.Values{}
	if params != nil && len(params.Types) > 0 {
		for _, t := range params.Types {
			query.Add("types", t)
		}
	} else {
		// Default to STANDARD type if not specified (required by API)
		query.Add("types", "STANDARD")
	}

	var balances []Balance
	path := fmt.Sprintf("/v4/profiles/%d/balances", profileID)
	err := s.client.Get(ctx, path, query, &balances)
	if err != nil {
		return nil, err
	}
	return balances, nil
}

// Get retrieves a specific balance.
// GET /v4/profiles/{profileId}/balances/{balanceId}
func (s *BalancesService) Get(ctx context.Context, profileID, balanceID int64) (*Balance, error) {
	var balance Balance
	path := fmt.Sprintf("/v4/profiles/%d/balances/%d", profileID, balanceID)
	err := s.client.Get(ctx, path, nil, &balance)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// GetByCurrency retrieves a balance by currency.
func (s *BalancesService) GetByCurrency(ctx context.Context, profileID int64, currency Currency) (*Balance, error) {
	balances, err := s.List(ctx, profileID, nil)
	if err != nil {
		return nil, err
	}

	for _, b := range balances {
		if b.Currency == currency {
			return &b, nil
		}
	}

	return nil, &APIError{StatusCode: 404, Message: "balance not found for currency"}
}

// Convert converts money between balances using a quote.
// POST /v2/profiles/{profileId}/balance-movements
func (s *BalancesService) Convert(ctx context.Context, profileID int64, quoteID string) error {
	req := ConvertBalanceRequest{QuoteID: quoteID}
	path := fmt.Sprintf("/v2/profiles/%d/balance-movements", profileID)
	return s.client.Post(ctx, path, req, nil)
}

// GetStatement retrieves the statement for a balance.
// GET /v1/profiles/{profileId}/balance-statements/{balanceId}/statement.json
func (s *BalancesService) GetStatement(ctx context.Context, profileID, balanceID int64, currency Currency, intervalStart, intervalEnd string) ([]BalanceStatement, error) {
	query := url.Values{}
	query.Set("currency", string(currency))
	query.Set("intervalStart", intervalStart)
	query.Set("intervalEnd", intervalEnd)

	var result struct {
		Transactions []BalanceStatement `json:"transactions"`
	}
	path := fmt.Sprintf("/v1/profiles/%d/balance-statements/%d/statement.json", profileID, balanceID)
	err := s.client.Get(ctx, path, query, &result)
	if err != nil {
		return nil, err
	}
	return result.Transactions, nil
}
