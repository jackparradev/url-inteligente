package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/jackparradev/url-inteligente/internal/service"
)

func TestHandler_ShortenURL_Success(t *testing.T) {
	// Setup
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	// Request body
	reqBody := ShortenRequest{
		URL: "https://www.google.com",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	// Create request
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	// Create recorder
	rr := httptest.NewRecorder()
	
	// Execute
	handler.ShortenURL(rr, req)
	
	// Verify response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	
	var response ShortenResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	
	if response.LongURL != reqBody.URL {
		t.Errorf("Expected %s, got %s", reqBody.URL, response.LongURL)
	}
	
	if !strings.HasPrefix(response.ShortURL, "http://localhost:8080/") {
		t.Errorf("Invalid short URL format: %s", response.ShortURL)
	}
}

func TestHandler_ShortenURL_InvalidMethod(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	rr := httptest.NewRecorder()
	
	handler.ShortenURL(rr, req)
	
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestHandler_ShortenURL_InvalidJSON(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	handler.ShortenURL(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestHandler_ShortenURL_InvalidURL(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	reqBody := ShortenRequest{
		URL: "invalid-url",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	handler.ShortenURL(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestHandler_RedirectURL_Success(t *testing.T) {
	// Setup
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	// Create a short URL first
	longURL := "https://www.google.com"
	shortCode, _ := shortener.CreateShortURL(longURL)
	
	// Create request
	req := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	rr := httptest.NewRecorder()
	
	// Execute
	handler.RedirectURL(rr, req)
	
	// Verify redirect
	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status 301, got %d", rr.Code)
	}
	
	location := rr.Header().Get("Location")
	if location != longURL {
		t.Errorf("Expected location %s, got %s", longURL, location)
	}
}

func TestHandler_RedirectURL_NotFound(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()
	
	handler.RedirectURL(rr, req)
	
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestHandler_RedirectURL_EmptyCode(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	
	handler.RedirectURL(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestHandler_RedirectURL_InvalidMethod(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	handler := NewHandler(shortener)
	
	req := httptest.NewRequest(http.MethodPost, "/abc123", nil)
	rr := httptest.NewRecorder()
	
	handler.RedirectURL(rr, req)
	
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestIsValidURL(t *testing.T) {
	validURLs := []string{
		"https://www.google.com",
		"http://example.com",
		"https://sub.domain.com/path",
		"http://localhost:8080/test",
	}
	
	invalidURLs := []string{
		"",
		"invalid",
		"ftp://example.com",
		"https://",
		"http://",
		"www.google.com",
		"https://noextension",
	}
	
	for _, url := range validURLs {
		if !isValidURL(url) {
			t.Errorf("Expected %s to be valid", url)
		}
	}
	
	for _, url := range invalidURLs {
		if isValidURL(url) {
			t.Errorf("Expected %s to be invalid", url)
		}
	}
}
