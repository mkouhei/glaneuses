package main

import (
	"net/http"

	"github.com/bitly/go-simplejson"
)

func (srv *service) restClient() (*simplejson.Json, error) {
	resp, err := http.Get(srv.uri)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	js, err := simplejson.NewFromReader(resp.Body)
	return js, nil
}
