package main

import (
	"fmt"
	"log"
	"time"

	"github.com/miguel-branco/goconfig"
)

func (conf *config) loadConfig(filepath string) {
	c, err := goconfig.ReadConfigFile(filepath)
	if err != nil {
		log.Fatal(c, err)
	}
	e, err := c.GetInt64("global", "expire")
	if err != nil {
		e = 30
	}
	conf.expire = time.Duration(e) * time.Minute

	for name, url := range srvMap {
		acct, err := c.GetString("account", name)
		var uri string
		if err != nil {
			acct = ""
			uri = ""
		} else if name == "pypi" {
			acct, err := c.GetString("account", name)
			if err != nil {
				acct = ""
				uri = ""
			}
			apikey, err := c.GetString("account", "libraries")
			if err != nil {
				apikey = ""
				uri = ""
			}
			uri = fmt.Sprintf(url, acct, apikey)
		} else {
			uri = fmt.Sprintf(url, acct)
		}
		srv := service{
			name,
			acct,
			uri,
		}
		conf.services = append(conf.services, srv)
	}
}
