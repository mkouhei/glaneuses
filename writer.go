package main

import (
	"log"
	"time"

	"github.com/bitly/go-simplejson"
)

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
