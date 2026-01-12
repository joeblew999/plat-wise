package wise

import (
	"context"
	"net/url"
	"strings"
)

// ExchangeRatesService handles exchange rate API calls.
type ExchangeRatesService struct {
	client *Client
}

// ExchangeRate represents an exchange rate.
type ExchangeRate struct {
	Rate   float64   `json:"rate"`
	Source Currency  `json:"source"`
	Target Currency  `json:"target"`
	Time   Timestamp `json:"time"`
}

// GetRateParams represents the parameters for getting exchange rates.
type GetRateParams struct {
	Source Currency
	Target Currency
	Time   string // ISO 8601 timestamp for historical rates
}

// Get retrieves the current exchange rate for a currency pair.
// GET /v1/rates
func (s *ExchangeRatesService) Get(ctx context.Context, source, target Currency) (*ExchangeRate, error) {
	rates, err := s.List(ctx, &GetRateParams{Source: source, Target: target})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, &APIError{StatusCode: 404, Message: "rate not found"}
	}
	return &rates[0], nil
}

// List retrieves exchange rates based on parameters.
// GET /v1/rates
func (s *ExchangeRatesService) List(ctx context.Context, params *GetRateParams) ([]ExchangeRate, error) {
	query := url.Values{}
	if params != nil {
		if params.Source != "" {
			query.Set("source", string(params.Source))
		}
		if params.Target != "" {
			query.Set("target", string(params.Target))
		}
		if params.Time != "" {
			query.Set("time", params.Time)
		}
	}

	var rates []ExchangeRate
	err := s.client.Get(ctx, "/v1/rates", query, &rates)
	if err != nil {
		return nil, err
	}
	return rates, nil
}

// GetHistorical retrieves a historical exchange rate at a specific time.
// GET /v1/rates?time=...
func (s *ExchangeRatesService) GetHistorical(ctx context.Context, source, target Currency, time string) (*ExchangeRate, error) {
	rates, err := s.List(ctx, &GetRateParams{Source: source, Target: target, Time: time})
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, &APIError{StatusCode: 404, Message: "rate not found"}
	}
	return &rates[0], nil
}

// HistoryParams represents parameters for getting rate history over a period.
type HistoryParams struct {
	Source Currency
	Target Currency
	From   string // ISO 8601 start timestamp
	To     string // ISO 8601 end timestamp
	Group  string // Interval: "day", "hour", or "minute"
}

// GetHistory retrieves exchange rate history over a period.
// GET /v1/rates?source=EUR&target=USD&from=...&to=...&group=day
func (s *ExchangeRatesService) GetHistory(ctx context.Context, params *HistoryParams) ([]ExchangeRate, error) {
	query := url.Values{}
	if params != nil {
		if params.Source != "" {
			query.Set("source", string(params.Source))
		}
		if params.Target != "" {
			query.Set("target", string(params.Target))
		}
		if params.From != "" {
			query.Set("from", params.From)
		}
		if params.To != "" {
			query.Set("to", params.To)
		}
		if params.Group != "" {
			query.Set("group", params.Group)
		}
	}

	var rates []ExchangeRate
	err := s.client.Get(ctx, "/v1/rates", query, &rates)
	if err != nil {
		return nil, err
	}
	return rates, nil
}

// GetMultiple retrieves rates for multiple currency pairs.
// Returns a map of "SOURCE-TARGET" -> rate
func (s *ExchangeRatesService) GetMultiple(ctx context.Context, pairs [][2]Currency) (map[string]float64, error) {
	// Build query with all pairs
	// Note: Wise API returns all available rates if no source/target specified
	rates, err := s.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Build lookup map
	rateMap := make(map[string]float64)
	for _, r := range rates {
		key := string(r.Source) + "-" + string(r.Target)
		rateMap[key] = r.Rate
	}

	// Filter for requested pairs if specified
	if len(pairs) > 0 {
		result := make(map[string]float64)
		for _, pair := range pairs {
			key := string(pair[0]) + "-" + string(pair[1])
			if rate, ok := rateMap[key]; ok {
				result[key] = rate
			}
		}
		return result, nil
	}

	return rateMap, nil
}

// ParseCurrencyPair parses a currency pair string like "USD-EUR" into source and target currencies.
func ParseCurrencyPair(pair string) (source, target Currency, ok bool) {
	parts := strings.Split(pair, "-")
	if len(parts) != 2 {
		return "", "", false
	}
	return Currency(parts[0]), Currency(parts[1]), true
}
