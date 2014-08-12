package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitly/go-simplejson"
)

func TestRestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		js := simplejson.New()
		js.Set("hello", "world")
		data, _ := js.EncodePretty()
		fmt.Fprintln(w, string(data))
	}))
	defer ts.Close()
	js := restClient(ts.URL)
	if js == nil {
		t.Errorf("%v", js)
	}
	if v, err := js.Get("hello").String(); err != nil {
		t.Errorf("%v, want: world\n", v)
	}
}
