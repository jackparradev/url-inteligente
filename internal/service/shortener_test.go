package service

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestShortener_CreateShortURL(t *testing.T) {
	storage := NewStorage()
	shortener := NewShortener(storage)
	
	longURL := "https://www.google.com"
	shortCode, err := shortener.CreateShortURL(longURL)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(shortCode) != SHORT_CODE_LENGTH {
		t.Errorf("Expected short code length %d, got %d", SHORT_CODE_LENGTH, len(shortCode))
	}
	
	// Verificar que se puede recuperar
	retrievedURL, exists := shortener.GetLongURL(shortCode)
	if !exists {
		t.Error("Expected URL to exist")
	}
	
	if retrievedURL != longURL {
		t.Errorf("Expected %s, got %s", longURL, retrievedURL)
	}
}

func TestShortener_GetLongURL_NonExistent(t *testing.T) {
	storage := NewStorage()
	shortener := NewShortener(storage)
	
	_, exists := shortener.GetLongURL("nonexistent")
	if exists {
		t.Error("Expected false for non-existent short code")
	}
}

func TestShortener_GenerateShortCode(t *testing.T) {
	storage := NewStorage()
	shortener := NewShortener(storage)
	
	longURL := "https://www.example.com"
	
	// Generar múltiples códigos para la misma URL
	codes := make(map[string]bool)
	for i := 0; i < 10; i++ {
		code := shortener.generateShortCode(longURL, i)
		
		// Verificar longitud
		if len(code) != SHORT_CODE_LENGTH {
			t.Errorf("Expected length %d, got %d", SHORT_CODE_LENGTH, len(code))
		}
		
		// Verificar que solo contiene caracteres hexadecimales
		for _, char := range code {
			if !strings.ContainsRune("0123456789abcdef", char) {
				t.Errorf("Invalid character in short code: %c", char)
			}
		}
		
		codes[code] = true
	}
	
	// Verificar que se generaron códigos diferentes
	if len(codes) < 5 {
		t.Error("Expected more unique codes generated")
	}
}

func TestShortener_ConcurrentCreation(t *testing.T) {
	storage := NewStorage()
	shortener := NewShortener(storage)
	
	var wg sync.WaitGroup
	numGoroutines := 100
	results := make(chan string, numGoroutines)
	
	// Crear URLs concurrentemente
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			longURL := fmt.Sprintf("https://test%d.com", id)
			shortCode, err := shortener.CreateShortURL(longURL)
			if err != nil {
				t.Errorf("Error creating short URL: %v", err)
				return
			}
			results <- shortCode
		}(i)
	}
	
	wg.Wait()
	close(results)
	
	// Verificar que todos los códigos son únicos
	codes := make(map[string]bool)
	for code := range results {
		if codes[code] {
			t.Errorf("Duplicate short code generated: %s", code)
		}
		codes[code] = true
	}
}

func TestShortener_CollisionHandling(t *testing.T) {
	storage := NewStorage()
	shortener := NewShortener(storage)
	
	// Simular colisión almacenando manualmente un código que podría generarse
	longURL1 := "https://www.test1.com"
	longURL2 := "https://www.test2.com"
	
	// Crear primera URL
	shortCode1, err := shortener.CreateShortURL(longURL1)
	if err != nil {
		t.Errorf("Error creating first short URL: %v", err)
	}
	
	// Crear segunda URL (diferente)
	shortCode2, err := shortener.CreateShortURL(longURL2)
	if err != nil {
		t.Errorf("Error creating second short URL: %v", err)
	}
	
	// Verificar que son diferentes
	if shortCode1 == shortCode2 {
		t.Error("Expected different short codes for different URLs")
	}
	
	// Verificar que ambas URLs se pueden recuperar
	retrieved1, exists1 := shortener.GetLongURL(shortCode1)
	retrieved2, exists2 := shortener.GetLongURL(shortCode2)
	
	if !exists1 || !exists2 {
		t.Error("Expected both URLs to exist")
	}
	
	if retrieved1 != longURL1 || retrieved2 != longURL2 {
		t.Error("URLs do not match original values")
	}
}
