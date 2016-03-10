package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
)

var server = &http.Server{Addr: ":8080"}

func getData(f func() ([]interface{}, error)) []interface{} {
	payload, err := f()
	if err != nil {
		log.Println(err)
	}
	return payload
}

func getJSON(f func() (*simplejson.Json, error)) *simplejson.Json {
	payload, err := f()
	if err != nil {
		log.Println(err)
	}
	return payload
}

func getPGP(f func() (pgp, error)) pgp {
	payload, err := f()
	if err != nil {
		log.Println(err)
	}
	return payload
}

func (conf *config) mergeJSON() ([]byte, error) {
	js := simplejson.New()
	for _, srv := range conf.services {
		switch {
		case srv.name == "debian":
			js.Set("deb", getData(srv.debPackages))
			srv.uri += "&format=json"
			if getJSON(srv.restClient) != nil {
				js.Set("udd", getJSON(srv.restClient).MustArray())
			}
		case srv.name == "github", srv.name == "bitbucket", srv.name == "rubygems":
			js.Set(srv.name, getJSON(srv.restClient))
		case srv.name == "pypi":
			js.Set(srv.name, getData(srv.pypiClient))
		case srv.name == "pgp":
			js.Set(srv.name, getPGP(srv.pgpData))
		}
	}

	js.Set("geneated_datetime", time.Now())
	data, err := js.EncodePretty()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (conf *config) app() {
	http.Handle("/", http.HandlerFunc(conf.serveJSON))
	if err := http.ListenAndServe(server.Addr, nil); err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}

func (conf *config) serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := conf.mergeJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(data))
}

func (conf *config) checkEtag(url string) string {
	resp, err := http.Head(url)
	if err != nil {
		return ""
	}
	if len(resp.Header["Etag"]) > 0 {
		return resp.Header["Etag"][0]
	}
	return fmt.Sprintf("expire-%d", int32(time.Now().Add(conf.expire).Unix()))
}
