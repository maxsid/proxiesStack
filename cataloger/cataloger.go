package cataloger

import (
	"crypto/tls"
	"github.com/maxsid/proxiesStack/config"
	"github.com/maxsid/proxiesStack/database"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func CategorizeProxies(proxiesChan chan string, done chan bool) {
	for {
		select {
		case proxyAddr := <-proxiesChan:
			client := GetClientWithProxy(proxyAddr)
			if CheckPageContentByClient(client, config.CheckPageAddress, config.CheckPattern) {
				_ = database.AddSetValue(database.WorkingSetKey, proxyAddr)
			} else {
				_ = database.AddSetValue(database.NotWorkingSetKey, proxyAddr)
			}
		case <-done:
			return
		}
	}
}

// Check a proxy connection and a content of the page by a client with proxy
func CheckPageContentByClient(client *http.Client, addr, regexPattern string) bool {
	resp, err := client.Get(addr)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Connection error: %s\n", err)
		return false
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll error: %v\n", err)
		return false
	}

	matched, err := regexp.Match(regexPattern, bodyBytes)
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
