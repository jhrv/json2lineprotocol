package transformer

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/astaxie/flatmap"
	"strings"
	"fmt"
	"time"
)

type Transformer struct {
	Request *http.Request
	Tags    map[string]string
}

func (t *Transformer) transform(url string) string {
	request := performRequest(t.Request)
	json, err := flatmap.Flatten(asJson(request))

	if (err != nil) {
		panic(err)
	}

	tags := createTagString(t.Tags)
	timestamp := time.Now().UnixNano()
	
	var lines []string
	for key, value := range json {
		line := fmt.Sprintf("%s%s value=%s %d", key, tags, value, timestamp)
		lines = append(lines, line)
	}

	result := strings.Join(lines, "\n")
	return result
}

func createTagString(tags map[string]string) string {
	var tagString string
	for key, value := range tags {
		tagString += "," + string(key) + "=" + string(value)
	}
	return tagString
}

func performRequest(req *http.Request) *http.Response {
	log.Printf("Calling %s \n", req.URL)
	client := &http.Client{}
	resp, err := client.Do(req)

	if (err != nil) {
		panic(err)
	}

	return resp
}

func asJson(resp *http.Response) map[string]interface{} {
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	d.UseNumber()

	var jsonData map[string]interface{}
	err := d.Decode(&jsonData)

	if (err != nil) {
		panic(err)
	}

	return jsonData
}
