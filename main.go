package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/kolo/xmlrpc"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	username = "mkouhei"
	email    = "mkouhei@palmtb.net"
	udd      = "http://udd.debian.org/dmd/"
	pypi     = "http://pypi.python.org/pypi"
	github   = "https://api.github.com/users"
)

type deb struct {
	Source string
	Url    string
}

type dl struct {
	LastDay   int `xmlrpc:"last_day"`
	LastMonth int `xmlrpc:"last_month"`
	LastWeek  int `xmlrpc:"last_week"`
}

func debPackages() []interface{} {
	doc, _ := goquery.NewDocument(udd + "?email1=" + email)
	cnt := doc.Find("h2#versions+table a").Length()
	debs := make([]interface{}, cnt)
	doc.Find("h2#versions+table a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			debs[i] = deb{s.Text(), url}
		}
	})
	return debs
}

func restClient(s string) *simplejson.Json {
	resp, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	js, err := simplejson.NewFromReader(resp.Body)
	return js
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
func pypiClient() []interface{} {
	client, _ := xmlrpc.NewClient(pypi, nil)
	defer client.Close()
	// PyPI user_packages()
	var result [][]string
	client.Call("user_packages", username, &result)
	pkgs := make([]interface{}, len(result))
	for i, v := range result {
		var ver []string
		client.Call("package_releases", v[1], &ver)
		meta := struct {
			Name       string `xmlrpc:"name"`
			Version    string `xmlrpc:"version"`
			PackageUrl string `xmlrpc:"package_url"`
			ReleaseUrl string `xmlrpc:"release_url"`
			Downloads  dl     `xmlrpc:"downloads"`
			Summary    string `xmlrpc:"summary"`
		}{}
		client.Call("release_data", []interface{}{v[1], ver[0]}, &meta)
		pkgs[i] = meta
	}
	return pkgs
}

func mergeJson() []byte {
	js := simplejson.New()
	js.Set("deb", debPackages())
	js.Set("udd", restClient(udd+"?email1="+email+"&format=json").MustArray())
	js.Set("github", restClient(github+"/"+username+"/events"))
	js.Set("pypi", pypiClient())
	data, _ := js.EncodePretty()
	return data
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, string(mergeJson()))
}

func main() {
	o := flag.String("o", "glaneuses.json", "Output file")
	flag.Parse()

	err := ioutil.WriteFile(*o, mergeJson(), 0644)
	if err != nil {
		panic(err)
	}
}
