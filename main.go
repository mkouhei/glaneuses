package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/kolo/xmlrpc"
	"log"
	"net/http"
)

var (
	username = "mkouhei"
	email    = "mkouhei@palmtb.net"
	udd      = "http://udd.debian.org/dmd/?format=json"
	pypi     = "http://pypi.python.org/pypi"
	mozilla  = "https://bugzilla.mozilla.org/xmlrpc.cgi"
	github   = "https://api.github.com/users"
)

func restClient(s string) []byte {
	resp, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	js, err := simplejson.NewFromReader(resp.Body)
	data, _ := js.EncodePretty()
	return data
}

func jsonRestClient(s string) interface{} {
	resp, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var data interface{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)
	return data
}

/*
confirm with cURL command and prepared XML file,
---
<?xml version="1.0"?>
<methodCall>
  <methodName>someMethod</methodName>
  <params>
    <param><string>someParam</string></param>
  </params>
</methodCall>
---
$ curl -H 'Content-Type: text/xml' -X POST --data @test.xml \
> http://example.org/xmlrpc/api

Response data types encoding rules is as follows;
https://github.com/kolo/xmlrpc#result-decoding
*/
func xmlRpcClient(s string) []string {
	client, _ := xmlrpc.NewClient(s, nil)
	defer client.Close()
	// PyPI user_packages()
	var result [][]string
	client.Call("user_packages", username, &result)
	pkgs := make([]string, len(result))
	for i, v := range result {
		pkgs[i] = v[1]
	}
	return pkgs
}

func assert(data interface{}) {
	switch data.(type) {
	case string:
		fmt.Print(data.(string))
	case float64:
		fmt.Print(data.(float64))
	case bool:
		fmt.Print(data.(bool))
	case nil:
		fmt.Print("null")
	case []interface{}:
		fmt.Print("[")
		for _, v := range data.([]interface{}) {
			assert(v)
			fmt.Print(" ")
		}
		fmt.Print("]")
	case map[string]interface{}:
		fmt.Print("{")
		for k, v := range data.(map[string]interface{}) {
			fmt.Print(k + ": ")
			assert(v)
			fmt.Print(" ")
		}
		fmt.Print("}")
	default:
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(restClient(udd+"&email1="+email)))
}

func main() {
	//uddResult := jsonRestClient(udd + "&email1=" + email)
	//restClient(udd + "&email1=" + email)
	//assert(uddResult)
	/*
		githubResult := jsonRestClient(github + "/" + username + "/events")
		assert(githubResult)
		pypiResult := xmlRpcClient(pypi)
		fmt.Println(pypiResult)
	*/
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
