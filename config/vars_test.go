package config

import (
	"testing"
)

func TestProcessHttpAddress(t *testing.T) {
	checkProcessedHttpAddress(t, "", "")
	checkProcessedHttpAddress(t, "google.com", "https://google.com")
	checkProcessedHttpAddress(t, "http://google.com", "http://google.com")
	checkProcessedHttpAddress(t, "https://google.com", "https://google.com")
}

func checkProcessedHttpAddress(t *testing.T, original, processed string) {
	processHttpAddress(&original)
	if original != processed {
		t.Errorf("address must be \"%s\", not \"%s\"", processed, original)
	}
}
