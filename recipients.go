package wise

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// RecipientsService handles recipient-related API calls.
type RecipientsService struct {
	client *Client
}

// Recipient represents a transfer recipient (account).
type Recipient struct {
	ID                int64                  `json:"id"`
	Profile           int64                  `json:"profile,omitempty"`
	AccountHolderName string                 `json:"accountHolderName"`
	Type              RecipientType          `json:"type"`
	Currency          Currency               `json:"currency"`
	Country           string                 `json:"country,omitempty"` // ISO 3166-1 alpha-2
	Active            bool                   `json:"active"`
	OwnedByCustomer   bool                   `json:"ownedByCustomer,omitempty"`
	Details           map[string]interface{} `json:"details"`
}

// CreateRecipientRequest represents the request to create a recipient.
type CreateRecipientRequest struct {
	Profile           int64                  `json:"profile"`
	AccountHolderName string                 `json:"accountHolderName"`
	Currency          Currency               `json:"currency"`
	Type              RecipientType          `json:"type"`
	Details           map[string]interface{} `json:"details"`
	OwnedByCustomer   bool                   `json:"ownedByCustomer,omitempty"`
}

// RecipientRequirements represents the requirements for creating a recipient.
type RecipientRequirements struct {
	Type   string                 `json:"type"`
	Title  string                 `json:"title,omitempty"`
	Fields []RecipientField       `json:"fields,omitempty"`
}

// RecipientField represents a field requirement for a recipient.
type RecipientField struct {
	Name          string                   `json:"name"`
	Group         []RecipientFieldGroup    `json:"group,omitempty"`
}

// RecipientFieldGroup represents a group of field validations.
type RecipientFieldGroup struct {
	Key               string            `json:"key"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	RefreshRequirementsOnChange bool    `json:"refreshRequirementsOnChange,omitempty"`
	Required          bool              `json:"required"`
	DisplayFormat     string            `json:"displayFormat,omitempty"`
	Example           string            `json:"example,omitempty"`
	MinLength         int               `json:"minLength,omitempty"`
	MaxLength         int               `json:"maxLength,omitempty"`
	ValidationRegexp  string            `json:"validationRegexp,omitempty"`
	ValidationAsync   interface{}       `json:"validationAsync,omitempty"`
	ValuesAllowed     []ValueAllowed    `json:"valuesAllowed,omitempty"`
}

// ValueAllowed represents an allowed value for a field.
type ValueAllowed struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ListParams represents the parameters for listing recipients.
type ListRecipientsParams struct {
	ProfileID int64
	Currency  Currency
	Limit     int
	Offset    int
}

// Create creates a new recipient.
// POST /v1/accounts
func (s *RecipientsService) Create(ctx context.Context, req *CreateRecipientRequest) (*Recipient, error) {
	var recipient Recipient
	err := s.client.Post(ctx, "/v1/accounts", req, &recipient)
	if err != nil {
		return nil, err
	}
	return &recipient, nil
}

// Get retrieves a recipient by ID.
// GET /v1/accounts/{accountId}
func (s *RecipientsService) Get(ctx context.Context, accountID int64) (*Recipient, error) {
	var recipient Recipient
	path := fmt.Sprintf("/v1/accounts/%d", accountID)
	err := s.client.Get(ctx, path, nil, &recipient)
	if err != nil {
		return nil, err
	}
	return &recipient, nil
}

// List returns all recipients for a profile.
// GET /v1/accounts
func (s *RecipientsService) List(ctx context.Context, params *ListRecipientsParams) ([]Recipient, error) {
	query := url.Values{}
	if params != nil {
		if params.ProfileID > 0 {
			query.Set("profile", strconv.FormatInt(params.ProfileID, 10))
		}
		if params.Currency != "" {
			query.Set("currency", string(params.Currency))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", strconv.Itoa(params.Offset))
		}
	}

	var recipients []Recipient
	err := s.client.Get(ctx, "/v1/accounts", query, &recipients)
	if err != nil {
		return nil, err
	}
	return recipients, nil
}

// Delete deletes a recipient by ID.
// DELETE /v1/accounts/{accountId}
func (s *RecipientsService) Delete(ctx context.Context, accountID int64) error {
	path := fmt.Sprintf("/v1/accounts/%d", accountID)
	return s.client.Delete(ctx, path, nil)
}

// GetRequirements returns the requirements for creating a recipient.
// GET /v1/account-requirements
func (s *RecipientsService) GetRequirements(ctx context.Context, quoteID string, currency Currency) ([]RecipientRequirements, error) {
	query := url.Values{}
	if quoteID != "" {
		query.Set("quoteId", quoteID)
	}
	if currency != "" {
		query.Set("targetCurrency", string(currency))
	}

	var requirements []RecipientRequirements
	err := s.client.Get(ctx, "/v1/account-requirements", query, &requirements)
	if err != nil {
		return nil, err
	}
	return requirements, nil
}
