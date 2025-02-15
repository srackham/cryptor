package mock

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	. "github.com/srackham/cryptor/internal/global"
)

func TestHttpGet(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		expectCode int
		expectBody string
	}{
		{"Valid BTC query", PRICE_QUERY + "BTC" + "USDT", http.StatusOK, `{"symbol":"BTCUSDT","price":"100000.00000000"}`},
		{"Valid ETH query", PRICE_QUERY + "ETH" + "USDT", http.StatusOK, `{"symbol":"ETHUSDT","price":"1000.00"}`},
		{"Valid USDC query", PRICE_QUERY + "USDC" + "USDT", http.StatusOK, `{"symbol":"USDCUSDT","price":"1.00"}`},
		{"Bad JSON response", PRICE_QUERY + "BAD_JSON" + "USDT", http.StatusOK, ``},
		{"Invalid symbol", PRICE_QUERY + "INVALID_SYMBOL" + "USDT", http.StatusBadRequest, ``},
		{"Valid exchange rate query", XRATES_QUERY + "1234", http.StatusOK, `{
  "rates": {
    "AUD": 1.6,
    "NZD": 1.5,
    "USD": 1
  }
}`},
		{"Unknown URL", "https://unknown.com", http.StatusNotFound, `not found`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := httpGet(tt.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.expectCode {
				t.Errorf("expected status %d, got %d", tt.expectCode, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			var expectedJSON, actualJSON interface{}
			json.Unmarshal([]byte(tt.expectBody), &expectedJSON)
			json.Unmarshal(body, &actualJSON)

			if expectedJSON == nil {
				if string(body) != tt.expectBody {
					t.Errorf("expected body %q, got %q", tt.expectBody, string(body))
				}
			} else if !equalJSON(expectedJSON, actualJSON) {
				t.Errorf("expected JSON %v, got %v", expectedJSON, actualJSON)
			}
		})
	}
}

// equalJSON compares JSON objects without regard to formatting
func equalJSON(a, b interface{}) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
