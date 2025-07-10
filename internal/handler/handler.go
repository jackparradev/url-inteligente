package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackparradev/url-inteligente/internal/service"
)

type Handler struct {
	shortener *service.Shortener
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(shortener *service.Shortener) *Handler {
	return &Handler{shortener: shortener}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validación básica de URL
	if !isValidURL(req.URL) {
		respondWithError(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Generar código corto
	shortCode, err := h.shortener.CreateShortURL(req.URL)
	if err != nil {
		respondWithError(w, "Error creating short URL", http.StatusInternalServerError)
		return
	}

	// Base URL hardcodeada según restricciones
	baseURL := "http://localhost:8080/"
	shortURL := baseURL + shortCode

	response := ShortenResponse{
		ShortURL: shortURL,
		LongURL:  req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer código corto de la URL
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	// Verificar que no esté vacío y no sea el endpoint /shorten
	if shortCode == "" || shortCode == "shorten" {
		respondWithError(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Buscar URL larga
	longURL, exists := h.shortener.GetLongURL(shortCode)
	if !exists {
		respondWithError(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redireccionar (301 - Moved Permanently)
	// Justificación: Para URLs acortadas, 301 es apropiado porque
	// el mapeo es permanente y permite caching del navegador
	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

// Función auxiliar para responder con errores
func respondWithError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Validación básica de URL (solo stdlib)
func isValidURL(url string) bool {
	// Validación básica sin librerías externas
	if len(url) < 8 {
		return false
	}

	// Debe empezar con http:// o https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// Debe tener al menos un punto después del protocolo
	if !strings.Contains(url[8:], ".") {
		return false
	}

	return true
}
