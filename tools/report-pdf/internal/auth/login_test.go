package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if req.Email != "user@test.com" {
			t.Errorf("Email = %q, want %q", req.Email, "user@test.com")
		}
		if req.Password != "secret" {
			t.Errorf("Password = %q, want %q", req.Password, "secret")
		}

		http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "tok123", Path: "/", HttpOnly: true})
		http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "ref456", Path: "/api/v1/refresh", HttpOnly: true})
		http.SetCookie(w, &http.Cookie{Name: "csrf_token", Value: "csrf789", Path: "/"})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"expires_in": 3600})
	}))
	defer server.Close()

	cookies, err := Login(server.URL, "user@test.com", "secret", "http://localhost:3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cookies) != 3 {
		t.Fatalf("got %d cookies, want 3", len(cookies))
	}

	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}

	if cookieMap["access_token"] != "tok123" {
		t.Errorf("access_token = %q", cookieMap["access_token"])
	}
	if cookieMap["refresh_token"] != "ref456" {
		t.Errorf("refresh_token = %q", cookieMap["refresh_token"])
	}
	if cookieMap["csrf_token"] != "csrf789" {
		t.Errorf("csrf_token = %q", cookieMap["csrf_token"])
	}

	// Verify domain is set to localhost for all cookies
	for _, c := range cookies {
		if *c.Domain != "localhost" {
			t.Errorf("cookie %q domain = %q, want %q", c.Name, *c.Domain, "localhost")
		}
	}
}

func TestLogin_ExtractsDomainFromBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "tok", Path: "/"})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cookies, err := Login(server.URL, "a@b.com", "pw", "https://app.example.com:3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}
	if *cookies[0].Domain != "app.example.com" {
		t.Errorf("domain = %q, want %q", *cookies[0].Domain, "app.example.com")
	}
}

func TestLogin_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := Login(server.URL, "a@b.com", "wrong", "http://localhost:3000")
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestLogin_NoCookies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"expires_in": 3600})
	}))
	defer server.Close()

	_, err := Login(server.URL, "a@b.com", "pw", "http://localhost:3000")
	if err == nil {
		t.Fatal("expected error when no cookies returned")
	}
}

func TestLogin_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := Login(server.URL, "a@b.com", "pw", "http://localhost:3000")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestLogin_ConnectionRefused(t *testing.T) {
	_, err := Login("http://localhost:1", "a@b.com", "pw", "http://localhost:3000")
	if err == nil {
		t.Fatal("expected error for connection refused")
	}
}

func TestLogin_PreservesCookiePaths(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "tok", Path: "/", HttpOnly: true})
		http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "ref", Path: "/api/v1/refresh", HttpOnly: true})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cookies, err := Login(server.URL, "a@b.com", "pw", "http://localhost:3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pathMap := make(map[string]string)
	for _, c := range cookies {
		pathMap[c.Name] = *c.Path
	}

	if pathMap["access_token"] != "/" {
		t.Errorf("access_token path = %q, want %q", pathMap["access_token"], "/")
	}
	if pathMap["refresh_token"] != "/api/v1/refresh" {
		t.Errorf("refresh_token path = %q, want %q", pathMap["refresh_token"], "/api/v1/refresh")
	}
}
