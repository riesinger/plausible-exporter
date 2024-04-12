package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riesinger/plausible-exporter/server"
)

func TestValidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testToken := "secret-token"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rr := httptest.NewRecorder()

	handlerWithMiddleware := server.BearerAuthMiddleware(handler, testToken)
	handlerWithMiddleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestInvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with an invalid token")
	})

	testToken := "secret-token"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	rr := httptest.NewRecorder()

	handlerWithMiddleware := server.BearerAuthMiddleware(handler, testToken)
	handlerWithMiddleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestMissingToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without a token")
	})

	testToken := "secret-token"

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handlerWithMiddleware := server.BearerAuthMiddleware(handler, testToken)
	handlerWithMiddleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}
