package util

import "net/url"

// Placeholder: Valida si una URL es válida
func IsValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}
