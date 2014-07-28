package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/kolo/xmlrpc"
	"github.com/miguel-branco/goconfig"
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

type deb struct {
	Source string
	Url    string
}

type pgp struct {
	Payload       string
	PublicKeyPath string
	VindexPath    string
}

type dl struct {
	LastDay   int `xmlrpc:"last_day"`
	LastMonth int `xmlrpc:"last_month"`
	LastWeek  int `xmlrpc:"last_week"`
}

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

func (a *Account) debPackages() ([]interface{}, error) {
	doc, err := goquery.NewDocument(udd + "?email1=" + a.DebianEmail)
	if err != nil {
		return nil, err
	}
	cnt := doc.Find("h2#versions+table a").Length()
	debs := make([]interface{}, cnt)
	doc.Find("h2#versions+table a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			debs[i] = deb{s.Text(), url}
		}
	})
	return debs, nil
}

func (a *Account) pgpData() (pgp, error) {
	keydata := &pgp{}
	doc, err := goquery.NewDocument(keyserver + a.KeyId)
	if err != nil {
		return *keydata, err
	}
	doc.Find("pre+hr+pre").Each(func(i int, s *goquery.Selection) {
		keydata.Payload = strings.Replace(
			strings.Replace(s.Text(), "@", " at ", -1),
			".", " dot ", -1)
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			url, exists := s.Attr("href")
			if exists {
				if i == 0 {
					keydata.PublicKeyPath = url
				} else if i == 1 {
					keydata.VindexPath = url
				}
			}
		})
	})
	return *keydata, nil
}

func restClient(s string) *simplejson.Json {
	resp, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	js, err := simplejson.NewFromReader(resp.Body)
	return js
}

/*
confirm with cURL command and prepared XML file,
---
<?xml version="1.0"?>
<methodCall>
  <methodName>someMethod</methodName>
  <params>
    <param><string>someParam</string></param>
  </params>
</methodCall>
---
$ curl -H 'Content-Type: text/xml' -X POST --data @test.xml \
> http://example.org/xmlrpc/api

Response data types encoding rules is as follows;
https://github.com/kolo/xmlrpc#result-decoding
*/
func (a *Account) pypiClient() ([]interface{}, error) {
	client, err := xmlrpc.NewClient(pypi, nil)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// PyPI user_packages()
	var result [][]string
	client.Call("user_packages", a.PypiUser, &result)
	pkgs := make([]interface{}, len(result))
	for i, v := range result {
		var ver []string
		client.Call("package_releases", v[1], &ver)
		meta := struct {
			Name       string `xmlrpc:"name"`
			Version    string `xmlrpc:"version"`
			PackageUrl string `xmlrpc:"package_url"`
			ReleaseUrl string `xmlrpc:"release_url"`
			Downloads  dl     `xmlrpc:"downloads"`
			Summary    string `xmlrpc:"summary"`
		}{}
		err := client.Call("release_data", []interface{}{v[1], ver[0]}, &meta)
		if err != nil {
			return nil, err
		}
		pkgs[i] = meta
	}
	return pkgs, nil
}

func (a *Account) mergeJson() ([]byte, error) {
	js := simplejson.New()
	dp, err := a.debPackages()
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
	pgpP, err := a.pgpData()
	js.Set("pgp", pgpP)
	data, err := js.EncodePretty()
	if err != nil {
		return nil, err
	}
	return data, nil
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
			log.Println("Gathering data and generate JSON.")
			data, err := a.mergeJson()
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile(*o, data, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
