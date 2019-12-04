package config

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Connection timeout
	TimeoutHTTP = 10
	// WebsiteAddress for grabbing
	CheckPageAddress = decode64("aHR0cHM6Ly9oaWRlbXkubmFtZS9lbi9kZW1vL3Jlc2V0Lw==")
	// Patter for match on the CheckPageAddress
	CheckPattern = decode64("PGZvcm0gYWN0aW9uPSJodHRwczovL2hpZGVteS5uYW1lL2VuL2RlbW8vc3VjY2Vzcy8iIG1ldGhvZD0icG9zdCIgaWQ9ImRmb3JtIg==")
	// Redis
	RedisHost string
	// Grab page
	GrabPageAddress = decode64("aHR0cHM6Ly9mcmVlLXByb3h5LWxpc3QubmV0")
	// Scan Interval
	ScanInterval = 300
	// Scan Off
	NoScan = false
)

// Parse environment
func ParseEnv() {
	setIntFromEnv(&TimeoutHTTP, "TIMEOUT")
	setStringFromEnv(&CheckPageAddress, "CHECK_ADDRESS")
	setStringFromEnv(&CheckPattern, "CHECK_PATTERN")
	setStringFromEnv(&GrabPageAddress, "GRAB_ADDRESS")
	setStringFromEnv(&RedisHost, "REDIS")
	setIntFromEnv(&ScanInterval, "SCAN_INTERVAL")
	setBoolFromEnv(&NoScan, "NO_SCAN")
	processParsed()
}

func setIntFromEnv(variable *int, envKey string) {
	envVar := getEnv(envKey)
	if envVar == "" {
		return
	}
	var err error
	*variable, err = strconv.Atoi(envVar)
	if err != nil {
		log.Panicf("Incorrect environment value %s. It must be int type, but it's %s.", envKey, envVar)
	}
}

func setBoolFromEnv(variable *bool, envKey string) {
	envVar := getEnv(envKey)
	if envVar == "" {
		return
	}
	var err error
	*variable, err = strconv.ParseBool(strings.ToLower(envVar))
	if err != nil {
		log.Panicf("Incorrect environment value %s. It must be int type, but it's %s.", envKey, envVar)
	}
}

func setStringFromEnv(variable *string, envKey string) {
	envVar := getEnv(envKey)
	if envVar == "" {
		return
	}
	*variable = envVar
}

func getEnv(envKey string) string {
	envVar := os.Getenv(envKey)
	_ = os.Unsetenv(envKey)
	return envVar
}

// Process all config variables
func processParsed() {
	processHttpAddress(&CheckPageAddress)
	processHttpAddress(&GrabPageAddress)
}

// Add a protocol name to the string start if it's not exist
func processHttpAddress(address *string) {
	if *address == "" {
		return
	}
	if hasProtocol, _ := regexp.MatchString(`^http(s)?://`, *address); !hasProtocol {
		*address = "https://" + *address
	}
}
