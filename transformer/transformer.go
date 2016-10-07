package transformer

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/astaxie/flatmap"
)

type Transformer struct {
	Request *http.Request
	Tags    map[string]string
}

//TODO tags sendes inn!

func (t *Transformer) transform(url string) string {
	jsonData := asJson(performRequest(t.Request))

	jsonData2, err := flatmap.Flatten(jsonData)

	if (err != nil) {
		panic(err)
	}

	var output string

	for key, value := range jsonData2 {
		log.Println(key, value)
		output += string(key) + ":" + string(value)
	}

	return output
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
