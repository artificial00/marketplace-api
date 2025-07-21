package utils

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// ValidateLogin проверяет корректность логина
func ValidateLogin(login string) bool {
	if len(login) < 3 || len(login) > 50 {
		return false
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", login)
	return matched
}

// ValidatePassword проверяет корректность пароля
func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}

// ValidateURL валидирует URL
func ValidateURL(rawURL string) bool {
	if rawURL == "" {
		return true // пустая строка допустима для nullable полей
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return false
	}

	if parsedURL.Host == "" {
		return false
	}

	return true
}
