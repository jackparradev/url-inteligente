package config

// Config contiene los parámetros de configuración de la aplicación.
type Config struct {
	// ServerPort es el puerto donde corre el servidor (incluye el prefijo ':').
	ServerPort string
	// BaseURL es la URL base usada para generar los enlaces cortos (sin barra final).
	BaseURL string
	// MaxRetry es el número máximo de intentos para generar un código único.
	MaxRetry int
	// ShortCodeLength es la longitud del código corto generado.
	ShortCodeLength int
}

// Get devuelve un puntero a Config con valores predefinidos.
func Get() *Config {
	return &Config{
		ServerPort:      ":8080",
		BaseURL:         "http://localhost:8080",
		MaxRetry:        5,
		ShortCodeLength: 6,
	}
}
