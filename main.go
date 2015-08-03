package main

import (
	"flag"
	"fmt"
	"log"
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

var srvMap map[string]string = map[string]string{
	"debian":    `https://udd.debian.org/dmd/?email1=%s`,
	"pypi":      "https://pypi.python.org/pypi",
	"rubygems":  `https://rubygems.org/api/v1/owners/%s/gems.json`,
	"github":    `https://api.github.com/users/%s/events`,
	"bitbucket": `https://bitbucket.org/api/1.0/users/%s/events`,
	"pgp":       `http://pgp.mit.edu/pks/lookup?op=index&fingerprint=on&search=%s`,
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

	conf := &config{}
	conf.loadConfig(*c)

	pollTicker := time.NewTicker(time.Duration(*p) * time.Minute)
	defer func() {
		pollTicker.Stop()
	}()
	for {
		select {
		case <-pollTicker.C:
			err := conf.writeJSON(*o)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
