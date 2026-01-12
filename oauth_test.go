package wise

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOAuthClient_AuthURL(t *testing.T) {
	client := NewOAuthClient(OAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Sandbox:      false,
	})

	url := client.AuthURL("test-state-123")

	// Check URL contains required params
	if url == "" {
		t.Error("AuthURL returned empty string")
	}
	if !contains(url, "client_id=test-client-id") {
		t.Error("AuthURL missing client_id")
	}
	if !contains(url, "redirect_uri=") {
		t.Error("AuthURL missing redirect_uri")
	}
	if !contains(url, "state=test-state-123") {
		t.Error("AuthURL missing state")
	}
	if !contains(url, "response_type=code") {
		t.Error("AuthURL missing response_type")
	}
	if !contains(url, ProductionAuthURL) {
		t.Errorf("AuthURL should use production URL, got: %s", url)
	}
}

func TestOAuthClient_AuthURL_Sandbox(t *testing.T) {
	client := NewOAuthClient(OAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Sandbox:      true,
	})

	url := client.AuthURL("state")

	if !contains(url, SandboxAuthURL) {
		t.Errorf("Sandbox AuthURL should use sandbox URL, got: %s", url)
	}
}

func TestOAuthClient_ExchangeCode(t *testing.T) {
	// Mock token server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Wrong content type: %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("Missing Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"access_token": "test-access-token",
			"token_type": "Bearer",
			"refresh_token": "test-refresh-token",
			"expires_in": 43200
		}`))
	}))
	defer server.Close()

	client := &OAuthClient{
		config: OAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/callback",
		},
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Override token URL for test
	token, err := client.tokenRequest(context.Background(), server.URL, map[string][]string{
		"grant_type": {"authorization_code"},
		"code":       {"test-auth-code"},
	})

	if err != nil {
		t.Fatalf("ExchangeCode failed: %v", err)
	}

	if token.AccessToken != "test-access-token" {
		t.Errorf("Wrong access token: %s", token.AccessToken)
	}
	if token.RefreshToken != "test-refresh-token" {
		t.Errorf("Wrong refresh token: %s", token.RefreshToken)
	}
	if token.ExpiresIn != 43200 {
		t.Errorf("Wrong expires_in: %d", token.ExpiresIn)
	}
}

func TestToken_IsExpired(t *testing.T) {
	// Token that expires in 1 hour - not expired
	token := &Token{
		AccessToken: "test",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	if token.IsExpired() {
		t.Error("Token with 1 hour remaining should not be expired")
	}

	// Token that expires in 2 minutes - should be considered expired (< 5 min buffer)
	token.ExpiresAt = time.Now().Add(2 * time.Minute)
	if !token.IsExpired() {
		t.Error("Token with 2 minutes remaining should be considered expired")
	}

	// Token that already expired
	token.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if !token.IsExpired() {
		t.Error("Token that expired 1 hour ago should be expired")
	}
}

func TestNewClientWithOAuth(t *testing.T) {
	token := &Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}

	client := NewClientWithOAuth(token, false)

	if client == nil {
		t.Fatal("NewClientWithOAuth returned nil")
	}
	if client.baseURL != ProductionBaseURL {
		t.Errorf("Expected production URL, got: %s", client.baseURL)
	}

	// Test sandbox
	client = NewClientWithOAuth(token, true)
	if client.baseURL != SandboxBaseURL {
		t.Errorf("Expected sandbox URL, got: %s", client.baseURL)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
