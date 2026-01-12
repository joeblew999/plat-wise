package wise

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuth endpoints
const (
	ProductionAuthURL  = "https://wise.com/oauth/authorize"
	ProductionTokenURL = "https://api.wise.com/oauth/token"
	SandboxAuthURL     = "https://sandbox.transferwise.tech/oauth/authorize"
	SandboxTokenURL    = "https://api.sandbox.transferwise.tech/oauth/token"
)

// OAuthConfig holds OAuth client credentials.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Sandbox      bool
	Scopes       []string
}

// Token represents an OAuth access token response.
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresIn    int       `json:"expires_in"`
	Scope        string    `json:"scope,omitempty"`
	ExpiresAt    time.Time `json:"-"`
}

// IsExpired returns true if the token is expired or about to expire.
func (t *Token) IsExpired() bool {
	// Consider expired if less than 5 minutes remaining
	return time.Now().Add(5 * time.Minute).After(t.ExpiresAt)
}

// OAuthClient handles OAuth authentication with Wise.
type OAuthClient struct {
	config     OAuthConfig
	httpClient *http.Client
}

// NewOAuthClient creates a new OAuth client.
func NewOAuthClient(config OAuthConfig) *OAuthClient {
	return &OAuthClient{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// AuthURL returns the authorization URL for user consent.
func (c *OAuthClient) AuthURL(state string) string {
	authURL := ProductionAuthURL
	if c.config.Sandbox {
		authURL = SandboxAuthURL
	}

	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("redirect_uri", c.config.RedirectURL)
	params.Set("response_type", "code")
	params.Set("state", state)
	if len(c.config.Scopes) > 0 {
		params.Set("scope", strings.Join(c.config.Scopes, " "))
	}

	return authURL + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for tokens.
func (c *OAuthClient) ExchangeCode(ctx context.Context, code string) (*Token, error) {
	tokenURL := ProductionTokenURL
	if c.config.Sandbox {
		tokenURL = SandboxTokenURL
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.config.RedirectURL)

	return c.tokenRequest(ctx, tokenURL, data)
}

// RefreshToken refreshes an expired access token.
func (c *OAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {
	tokenURL := ProductionTokenURL
	if c.config.Sandbox {
		tokenURL = SandboxTokenURL
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	return c.tokenRequest(ctx, tokenURL, data)
}

// ClientCredentials gets a token using client credentials (for server-to-server).
func (c *OAuthClient) ClientCredentials(ctx context.Context) (*Token, error) {
	tokenURL := ProductionTokenURL
	if c.config.Sandbox {
		tokenURL = SandboxTokenURL
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	return c.tokenRequest(ctx, tokenURL, data)
}

func (c *OAuthClient) tokenRequest(ctx context.Context, tokenURL string, data url.Values) (*Token, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Basic auth with client credentials
	auth := base64.StdEncoding.EncodeToString([]byte(c.config.ClientID + ":" + c.config.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("token request failed: %s - %s", resp.Status, string(body))
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}

	// Calculate expiration time
	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return &token, nil
}

// NewClientWithOAuth creates a new Wise API client with OAuth token.
func NewClientWithOAuth(token *Token, sandbox bool) *Client {
	opts := []ClientOption{}
	if sandbox {
		opts = append(opts, WithSandbox())
	}
	return NewClient(token.AccessToken, opts...)
}

// TokenManager handles automatic token refresh.
type TokenManager struct {
	oauth        *OAuthClient
	token        *Token
	onTokenRefresh func(*Token)
}

// NewTokenManager creates a token manager that auto-refreshes tokens.
func NewTokenManager(oauth *OAuthClient, initialToken *Token) *TokenManager {
	return &TokenManager{
		oauth: oauth,
		token: initialToken,
	}
}

// SetRefreshCallback sets a callback for when token is refreshed.
func (m *TokenManager) SetRefreshCallback(cb func(*Token)) {
	m.onTokenRefresh = cb
}

// GetToken returns a valid token, refreshing if needed.
func (m *TokenManager) GetToken(ctx context.Context) (*Token, error) {
	if m.token == nil {
		return nil, fmt.Errorf("no token available")
	}

	if !m.token.IsExpired() {
		return m.token, nil
	}

	if m.token.RefreshToken == "" {
		return nil, fmt.Errorf("token expired and no refresh token available")
	}

	newToken, err := m.oauth.RefreshToken(ctx, m.token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}

	m.token = newToken
	if m.onTokenRefresh != nil {
		m.onTokenRefresh(newToken)
	}

	return m.token, nil
}

// GetClient returns a Wise client with a valid token.
func (m *TokenManager) GetClient(ctx context.Context) (*Client, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	return NewClientWithOAuth(token, m.oauth.config.Sandbox), nil
}
