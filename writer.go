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
	js.Set("deb", dp)
	up := restClient(udd + "?email1=" + a.DebianEmail + "&format=json").MustArray()
	js.Set("udd", up)
	gp := restClient(github + a.GithubUser + "/events")
	js.Set("github", gp)
	bp := restClient(bitbucket + a.BitbucketUser + "/events")
	js.Set("bitbucket", bp)
	pp, err := a.pypiClient()
	js.Set("pypi", pp)
	rp := restClient(rubygems + a.GemsUser + "/gems.json")
	js.Set("rubygems", rp)
	pgpP, err := a.pgpData(keyserver)
	js.Set("pgp", pgpP)
	js.Set("geneated_datetime", time.Now())
	data, err := js.EncodePretty()
	if err != nil {
		return nil, err
	}
	return data, nil
}
