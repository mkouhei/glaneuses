package main

import (
	"fmt"
	"log"

	"github.com/miguel-branco/goconfig"
)

func (conf *config) loadConfig(filepath string) {
	c, err := goconfig.ReadConfigFile(filepath)
	if err != nil {
		log.Fatal(c, err)
	}

	for name, url := range srvMap {
		acct, err := c.GetString("account", name)
		var uri string
		if err != nil {
			acct = ""
			uri = ""
		} else if name == "pypi" {
			uri = url
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
