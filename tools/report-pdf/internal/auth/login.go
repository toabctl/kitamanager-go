package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/playwright-community/playwright-go"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates against the API and returns cookies suitable for Playwright.
func Login(apiURL, email, password, baseURL string) ([]playwright.OptionalCookie, error) {
	body, err := json.Marshal(loginRequest{Email: email, Password: password})
	if err != nil {
		return nil, fmt.Errorf("marshal login request: %w", err)
	}

	// Use a raw http.Client that does NOT follow redirects, so we capture Set-Cookie headers.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Post(apiURL+"/api/v1/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	domain := parsed.Hostname()

	var pwCookies []playwright.OptionalCookie
	for _, c := range resp.Cookies() {
		pwCookie := playwright.OptionalCookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   &domain,
			Path:     playwright.String(c.Path),
			HttpOnly: playwright.Bool(c.HttpOnly),
			Secure:   playwright.Bool(c.Secure),
			SameSite: playwright.SameSiteAttributeStrict,
		}
		pwCookies = append(pwCookies, pwCookie)
	}

	if len(pwCookies) == 0 {
		return nil, fmt.Errorf("no cookies received from login response")
	}

	return pwCookies, nil
}
