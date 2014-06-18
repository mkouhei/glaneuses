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
	github   = "https://api.github.com/users"
)

type userPackages struct {
	UserPackages string `xmlrpc:"user_packages"`
}

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

func xmlRpcClient(s string) {
	client, _ := xmlrpc.NewClient(s, nil)
	result := &userPackages{}
	client.Call("user_packages", "mkouhei", &result)
	fmt.Printf("%s\n", result.UserPackages)
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
	xmlRpcClient(pypi)
}
