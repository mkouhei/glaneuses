package main

import (
	"flag"
	"fmt"
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

type account struct {
	DebianEmail   string
	PypiUser      string
	GemsUser      string
	GithubUser    string
	BitbucketUser string
	KeyId         string
}

var version string
var showVersion = flag.Bool("version", false, "showVersion")

func main() {
	c := flag.String("c", "glaneuses.conf", "Configuration file")
	o := flag.String("o", "glaneuses.json", "Output file")
	p := flag.Int("p", 30, "Polling wait time (default: 30 (min))")
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", version)
		return
	}

	a := &account{}
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
