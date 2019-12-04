package grabber

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/maxsid/proxiesStack/config"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func GrabProxies(hostChan chan string, done chan bool) {
	resp, err := http.Get(config.GrabPageAddress)
	if err != nil {
		done <- true
		log.Fatalf("GrabProxies Error: %v\n", err)
	}
	node, err := html.Parse(resp.Body)
	if err != nil {
		done <- true
		log.Fatalf("GrabProxies Error: %v\n", err)
	}
	doc := goquery.NewDocumentFromNode(node)
	doc.Find("#proxylisttable tbody tr").Each(func(i int, selection *goquery.Selection) {
		host := selection.Find("td").First().Text()
		port := selection.Find("td").Eq(1).Text()
		isHttps := selection.Find("td.hx").Text() == "yes"
		if isHttps {
			hostChan <- fmt.Sprintf("https://%s:%s", host, port)
		} else {
			hostChan <- fmt.Sprintf("http://%s:%s", host, port)
		}
	})
	done <- true
}
