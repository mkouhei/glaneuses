package main

import (
	"net/http"

	"github.com/bitly/go-simplejson"
)

func restClient(s string) (*simplejson.Json, error) {
	resp, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	js, err := simplejson.NewFromReader(resp.Body)
	return js, nil
}
