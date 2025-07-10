package service

import (
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	urls map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		urls: make(map[string]string),
	}
}

func (s *Storage) Store(shortCode, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[shortCode] = longURL
}

func (s *Storage) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, exists := s.urls[shortCode]
	return longURL, exists
}

func (s *Storage) Exists(shortCode string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.urls[shortCode]
	return exists
}
