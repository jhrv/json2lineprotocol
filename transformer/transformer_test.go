package transformer

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
)

func TestSomething(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"foo": 69, "bar": {"pub": 96}}`)
	}))

	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	tags := map[string]string{"environment": "p", "application": "myapp" }

	transformer := &Transformer{req, tags}
	output := transformer.transform("http://nowhere.adeo.no")

	fmt.Println("output", output)

	expected := "foo value=69"

	if output != expected {
		t.Errorf("Mismatch between expected %s and actual %s", expected, output)
	}
}

/* cases

invalid json -> panic
multiple values -> works
tags

 */
