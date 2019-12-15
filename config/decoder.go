package config

import (
	"encoding/base64"
	"log"
)

// Decodes a base64 string line
func decode64(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal(err)
	}
	return string(decoded)
}
