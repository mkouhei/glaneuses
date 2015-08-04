package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type service struct {
	name    string
	account string
	uri     string
}

type config struct {
	services []service
}

var (
	srvMap map[string]string = map[string]string{
		"debian":    `https://udd.debian.org/dmd/?email1=%s`,
		"pypi":      "https://pypi.python.org/pypi",
		"rubygems":  `https://rubygems.org/api/v1/owners/%s/gems.json`,
		"github":    `https://api.github.com/users/%s/events`,
		"bitbucket": `https://bitbucket.org/api/1.0/users/%s/events`,
		"pgp":       `https://pgp.mit.edu/pks/lookup?op=index&fingerprint=on&search=%s`,
	}
	version     string
	showVersion = flag.Bool("version", false, "showVersion")
	server      = &http.Server{Addr: ":8080"}
)

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

func main() {
	c := flag.String("c", "glaneuses.conf", "Configuration file")
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", version)
		return
	}

	conf := &config{}
	conf.loadConfig(*c)
	conf.app()
}
