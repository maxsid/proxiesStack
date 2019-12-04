package config

import "testing"

func TestDecoder64(t *testing.T) {
	if decode64("aHR0cDovL2EuYi5jLmQucnU=") != "http://a.b.c.d.ru" {
		t.Error("config.decode64 is not work right!")
	}
}
