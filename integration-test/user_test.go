package integration_test

import (
	"context"
	"net/http"
	"testing"
)

// HTTP GET: /v1/user/profile.
func TestHTTPProfileV1(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token := registerAndLogin(t)

		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doAuthenticatedRequest(ctx, http.MethodGet, basePathV1+"/user/profile", http.NoBody, token)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		result := parseJSON[struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
		}](t, resp)

		if result.User.ID == "" {
			t.Error("Expected non-empty id")
		}
	})

	t.Run("no token", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, basePathV1+"/user/profile", http.NoBody)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})
}
