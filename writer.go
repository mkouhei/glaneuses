package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/bitly/go-simplejson"
)

func (a *Account) writeJson(o string) error {
	log.Println("Gathering data and generate JSON.")
	data, err := a.mergeJson()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(o, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *Account) mergeJson() ([]byte, error) {
	js := simplejson.New()
	dp, err := a.debPackages(udd)
	if err != nil {
		log.Println(err)
	}
	js.Set("deb", dp)

	up, err := restClient(udd + "?email1=" + a.DebianEmail + "&format=json")
	if err != nil {
		log.Println(err)
	}
	js.Set("udd", up.MustArray())

	gp, err := restClient(github + a.GithubUser + "/events")
	if err != nil {
		log.Println(err)
	}
	js.Set("github", gp)

	bp, err := restClient(bitbucket + a.BitbucketUser + "/events")
	if err != nil {
		log.Println(err)
	}
	js.Set("bitbucket", bp)

	pp, err := a.pypiClient()
	if err != nil {
		log.Println(err)
	}
	js.Set("pypi", pp)

	rp, err := restClient(rubygems + a.GemsUser + "/gems.json")
	if err != nil {
		log.Println(err)
	}
	js.Set("rubygems", rp)

	pgpP, err := a.pgpData(keyserver)
	if err != nil {
		log.Println(err)
	}
	js.Set("pgp", pgpP)

	js.Set("geneated_datetime", time.Now())
	data, err := js.EncodePretty()
	if err != nil {
		return nil, err
	}
	return data, nil
}
