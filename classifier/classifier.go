package classifier

import (
	"crypto/tls"
	"github.com/maxsid/proxiesStack/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func ClassifyProxies(proxiesChan <-chan string, workingChan, notWorkingChan chan<- string) {
	for proxyAddr := range proxiesChan {
		go func(proxyAddr string) {
			client := GetClientWithProxy(proxyAddr)
			body := getBodyViaClient(client, config.CheckPageAddress)
			if isMatchContent(body, config.CheckPagePattern) {
				log.Printf("%s will be added to the 'working' set", proxyAddr)
				workingChan <- proxyAddr
			} else {
				log.Printf("%s will be added to the 'not_working' set", proxyAddr)
				notWorkingChan <- proxyAddr
			}
		}(proxyAddr)
	}
}

// Check a proxy connection and a content of the page by a client with proxy
func getBodyViaClient(client *http.Client, addr string) []byte {
	resp, err := client.Get(addr)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Connection error: %s\n", err)
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll error: %v\n", err)
		return nil
	}
	return bodyBytes
}

// Checks match existing by Regex pattern
func isMatchContent(content []byte, regexPattern string) bool {
	matched, err := regexp.Match(regexPattern, content)
	if err != nil {
		log.Printf("Regexp Match error: %s\n", err)
		return false
	}
	return matched
}

func GetClientWithProxy(proxyAddr string) *http.Client {
	httpTransport := http.Transport{}
	if proxyAddr != "" {
		proxyUrl, err := url.Parse(proxyAddr)
		if err != nil {
			log.Println("Proxy address parse error: ", err)
		}

		httpTransport = http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &http.Client{
		Timeout:   time.Duration(config.TimeoutHTTP * int(time.Second)),
		Transport: &httpTransport,
	}
}
