package grabber

import (
	"github.com/maxsid/proxiesStack/config"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Gets GrabPage and looking for hosts on it
func GrabProxies(hostChan chan<- string) {
	resp, err := http.Get(config.GrabPageAddress)
	if err != nil {
		log.Printf("GrabProxies Error: %v\n", err)
	}
	pageBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GrabProxies Error: %v\n", err)
	}
	findHosts(pageBody, config.GrabPagePattern, hostChan)
	close(hostChan)
}

// Finds hosts on the page node by pattern
func findHosts(bodyPage []byte, pattern string, hostChan chan<- string) {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllSubmatch(bodyPage, -1)
	if matches == nil {
		log.Println("No matches")
		return
	}
	log.Printf("Found %d matches\n", len(matches))
	for _, match := range matches {
		matchGroups := make(map[string]string)
		for i, groupName := range re.SubexpNames() {
			if groupName != "" {
				matchGroups[groupName] = string(match[i])
			}
		}
		hostEnd := matchGroups["host"] + ":" + matchGroups["port"]
		if isHttps(matchGroups["https"]) {
			log.Printf("Add 'https://%s'", hostEnd)
			hostChan <- "https://" + hostEnd
		} else {
			log.Printf("Add 'http://%s'", hostEnd)
			hostChan <- "http://" + hostEnd
		}
	}
}

// Get true if value is Https
func isHttps(value string) bool {
	switch strings.ToLower(value) {
	case "yes", "ok", "https", "true":
		return true
	default:
		return false
	}
}
