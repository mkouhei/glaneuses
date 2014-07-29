package main

import (
	"flag"
	"log"
	"time"
)

const (
	udd       = "http://udd.debian.org/dmd/"
	pypi      = "http://pypi.python.org/pypi"
	rubygems  = "https://rubygems.org/api/v1/owners/"
	github    = "https://api.github.com/users/"
	bitbucket = "https://bitbucket.org/api/1.0/users/"
	keyserver = "http://pgp.mit.edu/pks/lookup?op=index&fingerprint=on&search="
)

type Account struct {
	DebianEmail   string
	PypiUser      string
	GemsUser      string
	GithubUser    string
	BitbucketUser string
	KeyId         string
}

func main() {
	c := flag.String("c", "glaneuses.conf", "Configuration file")
	o := flag.String("o", "glaneuses.json", "Output file")
	p := flag.Int("p", 30, "Polling wait time (default: 30 (min))")
	flag.Parse()

	a := &Account{}
	a.readConfig(*c)

	pollTicker := time.NewTicker(time.Duration(*p) * time.Minute)
	defer func() {
		pollTicker.Stop()
	}()
	for {
		select {
		case <-pollTicker.C:
			err := a.writeJson(*o)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
