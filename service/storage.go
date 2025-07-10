package service

import "sync"

var (
	urlMap  = make(map[string]string)
	urlLock = sync.RWMutex{}
)

// Placeholder: Guarda una URL en el mapa
func StoreURL(shortCode, longURL string) {
	urlLock.Lock()
	defer urlLock.Unlock()
	urlMap[shortCode] = longURL
}

// Placeholder: Recupera una URL larga
func GetURL(shortCode string) (string, bool) {
	urlLock.RLock()
	defer urlLock.RUnlock()
	longURL, ok := urlMap[shortCode]
	return longURL, ok
}
