package main

import (
	"log"
	"net/http"

	"github.com/jackparradev/url-inteligente/internal/handler"
	"github.com/jackparradev/url-inteligente/internal/service"
)

func main() {
	// Inicializar el storage
	storage := service.NewStorage()

	// Inicializar el servicio shortener
	shortener := service.NewShortener(storage)

	// Inicializar handlers
	h := handler.NewHandler(shortener)

	// Configurar rutas
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.ShortenURL)
	mux.HandleFunc("/", h.RedirectURL)

	// Puerto hardcodeado según restricciones
	port := ":8080"

	// Iniciar servidor
	log.Printf("Servidor iniciado en %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
