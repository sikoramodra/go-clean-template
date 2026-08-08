package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
)

// HTTP POST: /auth/signup.
func TestHTTPRegister(t *testing.T) {
	// Pre-register a user for the duplicate test case.
	name := sanitizeTestName(t)
	dupEmail := name + "_dup@test.com"

	resp := registerUser(t, dupEmail, testPassword)
	resp.Body.Close()

	// SuperTokens returns 200 on successful signup.
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("pre-register: expected 200, got %d", resp.StatusCode)
	}

	tests := []struct {
		description string
		email       string
		password    string
		// expectedHTTP is the HTTP status we accept.
		// SuperTokens often returns 200 even for business errors and puts
		// the real outcome in the JSON "status" field.
		expectedHTTP int
		// optional SuperTokens status string to assert when HTTP is 200.
		stStatus string
	}{
		{
			description:  "success",
			email:        name + "_ok@test.com",
			password:     testPassword,
			expectedHTTP: http.StatusOK,
			stStatus:     "OK",
		},
		{
			description:  "duplicate email",
			email:        dupEmail,
			password:     testPassword,
			expectedHTTP: http.StatusOK, // SuperTokens returns 200 + FIELD_ERROR
			stStatus:     "FIELD_ERROR",
		},
		{
			description:  "missing password",
			email:        name + "_nopw@test.com",
			password:     "",
			expectedHTTP: http.StatusOK, // validation error via formFields
			stStatus:     "FIELD_ERROR",
		},
		{
			description:  "short password",
			email:        name + "_short@test.com",
			password:     "ab", // SuperTokens default min length is 8
			expectedHTTP: http.StatusOK,
			stStatus:     "FIELD_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			resp := registerUser(t, tt.email, tt.password)
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedHTTP {
				t.Errorf("Expected HTTP status %d, got %d", tt.expectedHTTP, resp.StatusCode)
			}

			if tt.stStatus != "" {
				var body struct {
					Status string `json:"status"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode body: %v", err)
				}

				if body.Status != tt.stStatus {
					t.Errorf("Expected SuperTokens status %q, got %q", tt.stStatus, body.Status)
				}
			}

			// On success, tokens should be present in headers or cookies.
			if tt.stStatus == "OK" {
				if tok := extractAccessToken(resp); tok == "" {
					t.Error("expected access token in st-access-token header or sAccessToken cookie")
				}
			}
		})
	}
}

// HTTP POST: /auth/signin.
func TestHTTPLogin(t *testing.T) {
	email := sanitizeTestName(t) + "@test.com"
	password := testPassword

	// Register a user first.
	resp := registerUser(t, email, password)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("pre-register: expected 200, got %d", resp.StatusCode)
	}

	tests := []struct {
		description  string
		body         string
		expectedHTTP int
		stStatus     string
		checkToken   bool
	}{
		{
			description: "success",
			body: fmt.Sprintf(
				`{"formFields":[{"id":"email","value":%q},{"id":"password","value":%q}]}`,
				email, password,
			),
			expectedHTTP: http.StatusOK,
			stStatus:     "OK",
			checkToken:   true,
		},
		{
			description: "wrong password",
			body: fmt.Sprintf(
				`{"formFields":[{"id":"email","value":%q},{"id":"password","value":"wrongpass"}]}`,
				email,
			),
			expectedHTTP: http.StatusOK,
			stStatus:     "WRONG_CREDENTIALS_ERROR",
		},
		{
			description: "missing email",
			body: fmt.Sprintf(
				`{"formFields":[{"id":"password","value":%q}]}`,
				password,
			),
			expectedHTTP: http.StatusOK,
			stStatus:     "FIELD_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
			defer cancel()

			resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, httpURL+"/auth/signin", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedHTTP {
				t.Errorf("Expected HTTP status %d, got %d", tt.expectedHTTP, resp.StatusCode)
			}

			var body struct {
				Status string `json:"status"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode body: %v", err)
			}

			if body.Status != tt.stStatus {
				t.Errorf("Expected SuperTokens status %q, got %q", tt.stStatus, body.Status)
			}

			if tt.checkToken {
				// Body was already consumed; re-check via header/cookie on the response.
				if tok := extractAccessToken(resp); tok == "" {
					t.Error("Expected non-empty access token in st-access-token header or sAccessToken cookie")
				}
			}
		})
	}
}
