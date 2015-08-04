package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
)

var server = &http.Server{Addr: ":8080"}

func (conf *config) mergeJSON() ([]byte, error) {
	js := simplejson.New()
	for _, srv := range conf.services {
		switch {
		case srv.name == "debian":
			payload, err := srv.debPackages()
			if err != nil {
				log.Println(err)
			}
			js.Set("deb", payload)
			srv.uri += "&format=json"
			payload2, err := srv.restClient()
			if err != nil {
				log.Println(err)
			}
			js.Set("udd", payload2.MustArray())
		case (srv.name == "github" || srv.name == "bitbucket" || srv.name == "rubygems"):
			payload, err := srv.restClient()
			if err != nil {
				log.Println(err)
			}
			js.Set(srv.name, payload)
		case srv.name == "pypi":
			payload, err := srv.pypiClient()
			if err != nil {
				log.Println(err)
			}
			js.Set(srv.name, payload)
		case srv.name == "pgp":
			payload, err := srv.pgpData()
			if err != nil {
				log.Println(err)
			}
			js.Set(srv.name, payload)
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
