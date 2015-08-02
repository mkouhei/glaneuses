package main

import (
	"log"

	"github.com/miguel-branco/goconfig"
)

func (a *account) readConfig(p string) {
	c, err := goconfig.ReadConfigFile(p)
	if err != nil {
		log.Fatal(c, err)
	}
	a.DebianEmail, err = c.GetString("account", "debian")
	if err != nil {
		log.Println(err)
		a.DebianEmail = ""
	}
	a.PypiUser, err = c.GetString("account", "pypi")
	if err != nil {
		log.Println(err)
		a.PypiUser = ""
	}
	a.GemsUser, err = c.GetString("account", "rubygems")
	if err != nil {
		log.Println(err)
		a.GemsUser = ""
	}
	a.GithubUser, err = c.GetString("account", "github")
	if err != nil {
		log.Println(err)
		a.GithubUser = ""
	}
	a.BitbucketUser, err = c.GetString("account", "bitbucket")
	if err != nil {
		log.Println(err)
		a.BitbucketUser = ""
	}
	a.KeyID, err = c.GetString("account", "pgp")
	if err != nil {
		log.Println(err)
		a.KeyID = ""
	}
}
