package wise

import "time"

// Currency represents a currency code (ISO 4217).
type Currency string

// Common currency codes.
const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	AUD Currency = "AUD"
	CAD Currency = "CAD"
	CHF Currency = "CHF"
	CNY Currency = "CNY"
	INR Currency = "INR"
	SGD Currency = "SGD"
)

// Address represents a physical address.
type Address struct {
	FirstLine   string `json:"firstLine,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostCode    string `json:"postCode,omitempty"`
	Country     string `json:"country,omitempty"` // ISO 3166-1 alpha-2
	Occupation  string `json:"occupation,omitempty"`
}

// Money represents a monetary amount with currency.
type Money struct {
	Value    float64  `json:"value"`
	Currency Currency `json:"currency"`
}

// Timestamp is a time.Time that marshals to/from ISO 8601 format.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "null" || s == "" {
		return nil
	}

	// Try multiple formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05-0700",  // Wise format without colon in timezone
		"2006-01-02T15:04:05+0000",  // Wise UTC format
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	var err error
	for _, format := range formats {
		t.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON implements json.Marshaler.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

// TransferStatus represents the status of a transfer.
type TransferStatus string

const (
	TransferStatusIncomingPaymentWaiting  TransferStatus = "incoming_payment_waiting"
	TransferStatusIncomingPaymentInitiated TransferStatus = "incoming_payment_initiated"
	TransferStatusProcessing              TransferStatus = "processing"
	TransferStatusFundsConverted          TransferStatus = "funds_converted"
	TransferStatusOutgoingPaymentSent     TransferStatus = "outgoing_payment_sent"
	TransferStatusCancelled               TransferStatus = "cancelled"
	TransferStatusFundsRefunded           TransferStatus = "funds_refunded"
	TransferStatusBounced                 TransferStatus = "bounced_back"
)

// ProfileType represents the type of profile (personal or business).
type ProfileType string

const (
	ProfileTypePersonal ProfileType = "personal"
	ProfileTypeBusiness ProfileType = "business"
)

// RecipientType represents the type of recipient.
type RecipientType string

const (
	RecipientTypePerson  RecipientType = "person"
	RecipientTypeBusiness RecipientType = "business"
)
