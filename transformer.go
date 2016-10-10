package main

import (
	"net/http"
	"encoding/json"
	"github.com/astaxie/flatmap"
	"strings"
	"fmt"
	"time"
	"sort"
	"log"
)

type Transformer struct {
	Request *http.Request
	Tags    map[string]string
}

func (t *Transformer) Transform() string {
	request := performRequest(t.Request)
	json, err := flatmap.Flatten(asJson(request))
	log.Printf("flattened json to %s", json)

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

	log.Printf("transformed json to:\n%s\n-----", result)

	return result
}

func createTagString(tags map[string]string) string {
	var keys []string
	for k := range tags {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var tagString string
	for _, key := range keys {
		tagString += "," + string(key) + "=" + string(tags[key])
	}

	log.Printf("created tag string %s", tagString)
	return tagString
}

func performRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)

	if (err != nil) {
		panic(err)
	}

	log.Printf("got HTTP status %s from %s\n", resp.Status, req.URL)

	return resp
}

func asJson(resp *http.Response) map[string]interface{} {
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	d.UseNumber()

	var jsonData map[string]interface{}
	err := d.Decode(&jsonData)

	if (err != nil) {
		panic(fmt.Sprintf("unable to decode JSON. %s", err))
	}

	log.Printf("transformed response body to %s", jsonData)

	return jsonData
}
