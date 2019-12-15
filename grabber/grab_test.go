package grabber

import (
	"io/ioutil"
	"testing"
)

func TestFindHost(t *testing.T) {
	bodyPage, err := ioutil.ReadFile("../mockTest/grab-page/index.html")
	if err != nil {
		t.Fatal(err)
	}
	pattern := `<td>(?P<host>[\w\d\.].*?)</td>(?s:.*?)<td>(?P<port>\d+)(?s:.*?)` +
		`<td class="hm">(?P<https>yes|no)`

	hostChan := make(chan string)

	go findHosts(bodyPage, pattern, hostChan)
	mustBeResult := []string{"http://foo.com:1080", "http://bar.com:1080", "http://foo.net:1090", "http://bar.net:1090"}
	for _, host := range mustBeResult {
		foundHost := <-hostChan
		if host != foundHost {
			t.Fatalf("Must be %s, not %s", host, foundHost)
		}
	}

}

func TestIsHttps(t *testing.T) {
	values := []string{"https", "ok", "true", "yes", "http", "false", "no"}
	mustBe := []bool{true, true, true, true, false, false, false}
	for i, val := range values {
		if isHttps(val) != mustBe[i] {
			t.Errorf("isHttps('%s') has got the incorrect result (must be %v)", val, mustBe[i])
		}
	}
}
