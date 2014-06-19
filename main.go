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

type pkg struct {
	Role string
	Name string
}

type pkgs struct {
	Pkgs []pkg `xmlrpc:"array"`
}

func xmlRpcClient(s string) {
	client, _ := xmlrpc.NewClient(s, nil)
	defer client.Close()
	/*
			// Bugzill.version
			<?xml version="1.0" encoding="UTF-8"?>
			<methodResponse>
			<params>
			<param>
			<value>
			<struct>
			<member>
			<name>version</name>
			<value>
			<string>4.2.9+</string>
			</value></member>
			</struct>
			</value>
			</param>
			</params>
			</methodResponse>
		---
					var result = struct {
						Version string `xmlrpc:"version"`
					}{}
					client.Call("Bugzilla.version", nil, &result)
	*/

	/*
		//PyPI package_releases()
		---
		<?xml version='1.0'?>
		<methodResponse>
		<params>
		<param>
		<value><array><data>
		<value><string>0.3.4</string></value>
		</data></array></value>
		</param>
		</params>
		</methodResponse>
		---
			var result []string
			client.Call("package_releases", "shiori", &result)
	*/

	// PyPI user_packages()
	var result [][]string
	client.Call("user_packages", "mkouhei", &result)

	fmt.Println(result)
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
	client(udd + "&email1=" + email)
	client(github + "/" + username + "/events")
	//xmlRpcClient(mozilla)
	xmlRpcClient(pypi)
}
