package config

// Config contiene parámetros constantes definidos en código (no se cargan de archivo externo)
type Config struct {
	ServerPort      string // Puerto del servidor HTTP
	BaseURL         string // Dominio base para generar URLs cortas
	MaxRetry        int    // Máximo número de reintentos para evitar colisiones
	ShortCodeLength int    // Longitud del código corto
}

// Get retorna una configuración predefinida sin dependencias externas
func Get() Config {
	return Config{
		ServerPort:      ":8080",
		BaseURL:         "http://localhost:8080",
		MaxRetry:        5,
		ShortCodeLength: 6,
	}
}
