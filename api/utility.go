// Package api provides utility functions for the Go implementation of the API client.
package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomString generates a random string of a specified length.
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// PrintWelcomeMessage prints a welcome message to the console.
func PrintWelcomeMessage() {
	fmt.Println("Welcome to the Go implementation of the API client.")
}
