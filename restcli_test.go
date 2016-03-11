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

	srv := service{"test", "guest", ts.URL}
	js, err := srv.restClient()
	if err != nil {
		t.Errorf("%v", err)
	}

	if v, err := js.Get("hello").String(); err != nil {
		t.Errorf("%v, want: world\n", v)
	}
}
