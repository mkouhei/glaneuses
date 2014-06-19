package main

import (
	"encoding/json"
	"fmt"
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

func client(s string) {
	resp, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var data interface{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)
	assert(data)
	log.Println()
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
	client.Call("user_packages", "mkouhei", &result)
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

func main() {
	//client(udd + "&email1=" + email)
	//client(github + "/" + username + "/events")
	pkgs := xmlRpcClient(pypi)
	fmt.Println(pkgs)
}
