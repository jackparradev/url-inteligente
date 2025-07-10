package service

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

const (
	// Configuración hardcodeada según restricciones
	SHORT_CODE_LENGTH = 6
	MAX_ATTEMPTS      = 5
)

type Shortener struct {
	storage *Storage
}

func NewShortener(storage *Storage) *Shortener {
	// Inicializar seed para random
	rand.Seed(time.Now().UnixNano())
	return &Shortener{storage: storage}
}

func (s *Shortener) CreateShortURL(longURL string) (string, error) {
	// Intentar generar código único hasta MAX_ATTEMPTS veces
	for attempts := 0; attempts < MAX_ATTEMPTS; attempts++ {
		shortCode := s.generateShortCode(longURL, attempts)

		// Verificar si ya existe
		if !s.storage.Exists(shortCode) {
			s.storage.Store(shortCode, longURL)
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after %d attempts", MAX_ATTEMPTS)
}

func (s *Shortener) GetLongURL(shortCode string) (string, bool) {
	return s.storage.Get(shortCode)
}

func (s *Shortener) generateShortCode(longURL string, attempt int) string {
	// Combinar URL + timestamp + attempt para evitar colisiones
	// Solo usando librerías estándar
	input := fmt.Sprintf("%s%d%d", longURL, time.Now().UnixNano(), attempt)

	// Hash SHA1 (librería estándar)
	hash := sha1.Sum([]byte(input))

	// Convertir a string hexadecimal y tomar los primeros SHORT_CODE_LENGTH caracteres
	shortCode := fmt.Sprintf("%x", hash)[:SHORT_CODE_LENGTH]

	// Si es un reintento, agregar aleatoriedad extra
	if attempt > 0 {
		// Agregar 2 caracteres aleatorios al final
		randomSuffix := rand.Intn(256) // 0-255
		shortCode = fmt.Sprintf("%s%02x", shortCode[:SHORT_CODE_LENGTH-2], randomSuffix)
	}

	return shortCode
}
