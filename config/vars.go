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
	CheckPagePattern = decode64("PGZvcm0gYWN0aW9uPSJodHRwczovL2hpZGVteS5uYW1lL2VuL2RlbW8vc3VjY2Vzcy8i" +
		"IG1ldGhvZD0icG9zdCIgaWQ9ImRmb3JtIg==")
	// Redis
	RedisHost string
	// Grab page
	GrabPageAddress = decode64("aHR0cHM6Ly9mcmVlLXByb3h5LWxpc3QubmV0")
	GrabPagePattern = decode64("PHRyPjx0ZD4oP1A8aG9zdD4oPzpcd3xcZHxcLikrKTwvdGQ+KD9zOi4qPyk8dGQ+KD9QP" +
		"HBvcnQ+XGQrKTwvdGQ+KD9zOi4qPyk8dGQgY2xhc3M9J2h4Jz4oP1A8aHR0cHM+KD86eWVzfG5vKSkoP3M6Lio/KTwvdHI+")
	// Scan Interval in seconds
	ScanInterval = 300
	// Scan Off
	NoScan     = false
	ScanStatus = "off"
)

// Fills the variables from system Environment
func ParseEnv() {
	setIntFromEnv(&TimeoutHTTP, "TIMEOUT")
	setStringFromEnv(&CheckPageAddress, "CHECK_ADDRESS")
	setStringFromEnv(&CheckPagePattern, "CHECK_PATTERN")
	setStringFromEnv(&GrabPageAddress, "GRAB_ADDRESS")
	setStringFromEnv(&GrabPagePattern, "GRAB_PATTERN")
	setStringFromEnv(&RedisHost, "REDIS_HOST")
	setIntFromEnv(&ScanInterval, "SCAN_INTERVAL")
	setBoolFromEnv(&NoScan, "NO_SCAN")
	processParsed()
}

// Gets integer variable, panic error or nothing if variable is not exists
func setIntFromEnv(variable *int, envKey string) {
	envVar := getEnvVar(envKey)
	if envVar == "" {
		return
	}
	var err error
	*variable, err = strconv.Atoi(envVar)
	if err != nil {
		log.Panicf("Incorrect environment value %s. It must be int type, but it's %s.", envKey, envVar)
	}
}

// Gets boolean variable, panic error or nothing if variable is not exists
func setBoolFromEnv(variable *bool, envKey string) {
	envVar := getEnvVar(envKey)
	if envVar == "" {
		return
	}
	var err error
	*variable, err = strconv.ParseBool(strings.ToLower(envVar))
	if err != nil {
		log.Panicf("Incorrect environment value %s. It must be int type, but it's %s.", envKey, envVar)
	}
}

// Gets string variable or nothing if variable is not exists
func setStringFromEnv(variable *string, envKey string) {
	envVar := getEnvVar(envKey)
	if envVar == "" {
		return
	}
	*variable = envVar
}

// Gets variable from OS environment and cleans
func getEnvVar(envKey string) string {
	envVar := os.Getenv(envKey)
	_ = os.Unsetenv(envKey)
	return envVar
}

// Process all config variables
func processParsed() {
	processHttpAddress(&CheckPageAddress)
	processHttpAddress(&GrabPageAddress)
}

// Adds a protocol name to the string start if it's not exist
func processHttpAddress(address *string) {
	if *address == "" {
		return
	}
	if hasProtocol, _ := regexp.MatchString(`^http(s)?://`, *address); !hasProtocol {
		*address = "https://" + *address
	}
}
