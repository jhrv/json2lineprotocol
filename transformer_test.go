package transformer

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"strings"
)

func TestHappyPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"foo": 69}`)
	}))

	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	tags := map[string]string{"environment": "p", "application": "myapp" }

	transformer := Transformer{req, tags}
	result := transformer.Transform()

	expected := "foo,application=myapp,environment=p value=69"

	if !strings.HasPrefix(result, expected) {
		t.Errorf("Mismatch between expected and actual. Result %s does not start with expected %s", result, expected)
	}
}

func TestTagSorting(t *testing.T) {
	tags := map[string]string{"c": "1", "b": "2", "a": "3"}
	result := createTagString(tags)
	expected := ",a=3,b=2,c=1"
	if result != expected {
		t.Errorf("Mismatch between expected %s and actual %s", expected, result)
	}
}

func TestInvalidJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<xml />`)
	}))

	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	transformer := Transformer{req, nil}

	defer func(){
		err := recover()
		if err == nil {
			t.Errorf("Did not panic when endpoint supplies xml")
		}
	}()

	transformer.Transform()
}