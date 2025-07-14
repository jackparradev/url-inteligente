package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"github.com/jackparradev/url-inteligente/internal/handler"
	"github.com/jackparradev/url-inteligente/internal/service"
	"fmt"
)

func TestIntegration_FullWorkflow(t *testing.T) {
	// Setup completo
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	h := handler.NewHandler(shortener)
	
	// Configurar router
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.ShortenURL)
	mux.HandleFunc("/", h.RedirectURL)
	
	// Test 1: Crear URL corta
	reqBody := handler.ShortenRequest{
		URL: "https://www.example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	
	var response handler.ShortenResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	
	// Extraer código corto
	parts := strings.Split(response.ShortURL, "/")
	shortCode := parts[len(parts)-1]
	
	// Test 2: Usar URL corta para redirección
	redirectReq := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	redirectRR := httptest.NewRecorder()
	
	mux.ServeHTTP(redirectRR, redirectReq)
	
	if redirectRR.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status 301, got %d", redirectRR.Code)
	}
	
	location := redirectRR.Header().Get("Location")
	if location != reqBody.URL {
		t.Errorf("Expected location %s, got %s", reqBody.URL, location)
	}
}

func TestIntegration_ConcurrentRequests(t *testing.T) {
	// Setup
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	h := handler.NewHandler(shortener)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.ShortenURL)
	mux.HandleFunc("/", h.RedirectURL)
	
	var wg sync.WaitGroup
	numRequests := 100
	
	// Crear múltiples URLs concurrentemente
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func(id int) {
			defer wg.Done()
			
			reqBody := handler.ShortenRequest{
				URL: fmt.Sprintf("https://www.test%d.com", id),
			}
			jsonBody, _ := json.Marshal(reqBody)
			
			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			
			mux.ServeHTTP(rr, req)
			
			if rr.Code != http.StatusOK {
				t.Errorf("Request %d failed with status %d", id, rr.Code)
			}
		}(i)
	}
	
	wg.Wait()
}

func TestIntegration_MultipleURLsSameEndpoint(t *testing.T) {
	storage := service.NewStorage()
	shortener := service.NewShortener(storage)
	h := handler.NewHandler(shortener)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.ShortenURL)
	mux.HandleFunc("/", h.RedirectURL)
	
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
	}
	
	shortCodes := make([]string, len(urls))
	
	// Crear múltiples URLs
	for i, url := range urls {
		reqBody := handler.ShortenRequest{URL: url}
		jsonBody, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		
		mux.ServeHTTP(rr, req)
		
		var response handler.ShortenResponse
		json.NewDecoder(rr.Body).Decode(&response)
		
		parts := strings.Split(response.ShortURL, "/")
		shortCodes[i] = parts[len(parts)-1]
	}
	
	// Verificar que todos los códigos son únicos
	for i := 0; i < len(shortCodes); i++ {
		for j := i + 1; j < len(shortCodes); j++ {
			if shortCodes[i] == shortCodes[j] {
				t.Errorf("Duplicate short codes: %s", shortCodes[i])
			}
		}
	}
	
	// Verificar que todas las redirecciones funcionan
	for i, shortCode := range shortCodes {
		req := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
		rr := httptest.NewRecorder()
		
		mux.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusMovedPermanently {
			t.Errorf("Expected status 301 for %s, got %d", shortCode, rr.Code)
		}
		
		location := rr.Header().Get("Location")
		if location != urls[i] {
			t.Errorf("Expected location %s, got %s", urls[i], location)
		}
	}
}