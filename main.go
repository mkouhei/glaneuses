package main

import (
	"flag"
	"fmt"
	"time"
)

type service struct {
	name    string
	account string
	uri     string
}

type config struct {
	services []service
	expire   time.Duration
}

var (
	srvMap map[string]string = map[string]string{
		"debian":    `https://udd.debian.org/dmd/?email1=%s`,
		"pypi":      "https://pypi.python.org/pypi",
		"rubygems":  `https://rubygems.org/api/v1/owners/%s/gems.json`,
		"hackage":   `http://hackage.haskell.org/user/%s.json`,
		"github":    `https://api.github.com/users/%s/events`,
		"bitbucket": `https://api.bitbucket.org/2.0/repositories/%s`,
		"pgp":       `https://pgp.mit.edu/pks/lookup?op=index&fingerprint=on&search=%s`,
	}
	ver         string
	showVersion = flag.Bool("version", false, "showVersion")
	ignoreUids  []string
)

func main() {
	c := flag.String("c", "glaneuses.conf", "Configuration file")
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", ver)
		return
	}

	conf := &config{}
	conf.loadConfig(*c)
	conf.app()
}
