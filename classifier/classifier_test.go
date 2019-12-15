package classifier

import (
	"io/ioutil"
	"testing"
)

func TestIsMatchContent(t *testing.T) {
	bodyPage, err := ioutil.ReadFile("../mockTest/check-page/index.html")
	if err != nil {
		t.Fatal(err)
	}
	pattern := "<form.+?id=\"get-code\">"
	if !isMatchContent(bodyPage, pattern) {
		t.Error("check-page/index.html must be True")
	}

	if isMatchContent([]byte{}, `\w`) {
		t.Error("Empty content must return False")
	}
	if !isMatchContent([]byte("hello"), `h[a-z]+?o`) {
		t.Error("'Hello' content must return True")
	}
}
