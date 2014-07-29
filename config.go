package main

import (
	"log"

	"github.com/miguel-branco/goconfig"
)

func (a *Account) readConfig(p string) {
	c, err := goconfig.ReadConfigFile(p)
	if err != nil {
		log.Fatal(c, err)
	}
	a.DebianEmail, err = c.GetString("debian", "email")
	if err != nil {
		log.Println(err)
		a.DebianEmail = ""
	}
	a.PypiUser, err = c.GetString("pypi", "username")
	if err != nil {
		log.Println(err)
		a.PypiUser = ""
	}
	a.GemsUser, err = c.GetString("rubygems", "username")
	if err != nil {
		log.Println(err)
		a.GemsUser = ""
	}
	a.GithubUser, err = c.GetString("github", "username")
	if err != nil {
		log.Println(err)
		a.GithubUser = ""
	}
	a.BitbucketUser, err = c.GetString("bitbucket", "username")
	if err != nil {
		log.Println(err)
		a.BitbucketUser = ""
	}
	a.KeyId, err = c.GetString("pgp", "keyid")
	if err != nil {
		log.Println(err)
		a.KeyId = ""
	}
}
