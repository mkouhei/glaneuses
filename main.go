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
	services   []service
	expire     time.Duration
	ignoreUIDS []string
}

var (
	srvMap map[string]string = map[string]string{
		"debian":    `https://udd.debian.org/dmd/?email1=%s`,
		"pypi":      `https://libraries.io/api/search?platforms=Pypi&q=%s&api_key=%s`,
		"rubygems":  `https://rubygems.org/api/v1/owners/%s/gems.json`,
		"hackage":   `http://hackage.haskell.org/user/%s.json`,
		"github":    `https://api.github.com/users/%s/events`,
		"bitbucket": `https://api.bitbucket.org/2.0/repositories/%s`,
	}
	ver         string
	showVersion = flag.Bool("version", false, "showVersion")
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
	data, err := conf.mergeJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
