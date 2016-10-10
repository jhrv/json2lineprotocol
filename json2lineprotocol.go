package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	t "./transformer"
)

var (
	endpoint = flag.String("endpoint", "", "")
	tags = flag.String("tags", "", "")
	debug = flag.Bool("debug", false, "")
)

var usage = `Usage: json2lineprotocol [options...]

Options:

  -endpoint     "endpoint to extract JSON data from"
  -tags  	"tags on the format "key=value,key=value..."
  -debug 	"if set to true, execution outputs detailed descriptions of whats going on"
`

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	if !*debug {
		log.SetOutput(ioutil.Discard)
	}

	req, _ := http.NewRequest("GET", *endpoint, nil)

	transformer := t.Transformer{req, mapifyTagString(*tags)}
	output := transformer.Transform()
	fmt.Println(output)
}

func mapifyTagString(tagString string) map[string]string {
	tags := make(map[string]string)

	for _, tag := range strings.Split(tagString, ",") {
		splitted := strings.Split(tag, "=")
		key := splitted[0]
		value := splitted[1]
		tags[key] = value
	}

	return tags
}

