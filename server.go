package main

import (
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

func (conf *config) mergeJSON() ([]byte, error) {
	js := simplejson.New()
	for _, srv := range conf.services {
		if srv.name == "debian" {
			js.Set("deb", getData(srv.debPackages))
			srv.uri += "&format=json"
			if getJSON(srv.restClient) != nil {
				js.Set("udd", getJSON(srv.restClient).MustArray())
			}
		} else {
			js.Set(srv.name, getJSON(srv.restClient))
		}
	}

	js.Set("geneated_datetime", time.Now())
	data, err := js.EncodePretty()
	if err != nil {
		return nil, err
	}
	return data, nil
}
