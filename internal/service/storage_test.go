package service

import (
	"sync"
	"testing"
	"fmt"
)

func TestStorage_Store(t *testing.T) {
	storage := NewStorage()
	
	shortCode := "abc123"
	longURL := "https://www.google.com"
	
	storage.Store(shortCode, longURL)
	
	// Verificar que se almacenó correctamente
	retrievedURL, exists := storage.Get(shortCode)
	if !exists {
		t.Error("Expected URL to exist in storage")
	}
	
	if retrievedURL != longURL {
		t.Errorf("Expected %s, got %s", longURL, retrievedURL)
	}
}

func TestStorage_GetNonExistent(t *testing.T) {
	storage := NewStorage()
	
	_, exists := storage.Get("nonexistent")
	if exists {
		t.Error("Expected false for non-existent key")
	}
}

func TestStorage_Exists(t *testing.T) {
	storage := NewStorage()
	
	// Verificar que no existe inicialmente
	if storage.Exists("test123") {
		t.Error("Expected false for non-existent key")
	}
	
	// Almacenar y verificar que existe
	storage.Store("test123", "https://test.com")
	if !storage.Exists("test123") {
		t.Error("Expected true for existing key")
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	storage := NewStorage()
	
	var wg sync.WaitGroup
	numGoroutines := 100
	
	// Prueba de escritura concurrente
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			shortCode := fmt.Sprintf("code%d", id)
			longURL := fmt.Sprintf("https://test%d.com", id)
			storage.Store(shortCode, longURL)
		}(i)
	}
	
	wg.Wait()
	
	// Verificar que todas las URLs se almacenaron
	for i := 0; i < numGoroutines; i++ {
		shortCode := fmt.Sprintf("code%d", i)
		expectedURL := fmt.Sprintf("https://test%d.com", i)
		
		retrievedURL, exists := storage.Get(shortCode)
		if !exists {
			t.Errorf("Expected URL %s to exist", shortCode)
		}
		if retrievedURL != expectedURL {
			t.Errorf("Expected %s, got %s", expectedURL, retrievedURL)
		}
	}
}

func TestStorage_ConcurrentReadWrite(t *testing.T) {
	storage := NewStorage()
	
	// Almacenar algunos datos iniciales
	for i := 0; i < 10; i++ {
		shortCode := fmt.Sprintf("initial%d", i)
		longURL := fmt.Sprintf("https://initial%d.com", i)
		storage.Store(shortCode, longURL)
	}
	
	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 50
	
	// Lectores concurrentes
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				shortCode := fmt.Sprintf("initial%d", j)
				storage.Get(shortCode)
				storage.Exists(shortCode)
			}
		}(i)
	}
	
	// Escritores concurrentes
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			shortCode := fmt.Sprintf("writer%d", id)
			longURL := fmt.Sprintf("https://writer%d.com", id)
			storage.Store(shortCode, longURL)
		}(i)
	}
	
	wg.Wait()
	
	// Verificar integridad de datos
	for i := 0; i < 10; i++ {
		shortCode := fmt.Sprintf("initial%d", i)
		expectedURL := fmt.Sprintf("https://initial%d.com", i)
		
		retrievedURL, exists := storage.Get(shortCode)
		if !exists || retrievedURL != expectedURL {
			t.Errorf("Data integrity compromised for %s", shortCode)
		}
	}
}