package main

import (
	"flag"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/kolo/xmlrpc"
	"github.com/msbranco/goconfig"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	udd       = "http://udd.debian.org/dmd/"
	pypi      = "http://pypi.python.org/pypi"
	github    = "https://api.github.com/users"
	keyserver = "http://pgp.mit.edu/pks/lookup?op=index&fingerprint=on&search="
)

type Account struct {
	DebianEmail string
	PypiUser    string
	GithubUser  string
	KeyId       string
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

func readConfig(p string) Account {
	c, err := goconfig.ReadConfigFile(p)
	if err != nil {
		log.Fatal(c, err)
	}
	debianEmail, err := c.GetString("debian", "email")
	if err != nil {
		log.Fatal(err)
	}
	pypiUser, err := c.GetString("pypi", "username")
	if err != nil {
		log.Fatal(err)
	}
	githubUser, err := c.GetString("github", "username")
	if err != nil {
		log.Fatal(err)
	}
	keyId, err := c.GetString("pgp", "keyid")
	if err != nil {
		log.Fatal(err)
	}
	var acct Account = Account{debianEmail, pypiUser, githubUser, keyId}
	return acct
}

func debPackages(email string) []interface{} {
	doc, _ := goquery.NewDocument(udd + "?email1=" + email)
	cnt := doc.Find("h2#versions+table a").Length()
	debs := make([]interface{}, cnt)
	doc.Find("h2#versions+table a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			debs[i] = deb{s.Text(), url}
		}
	})
	return debs
}

func pgpData(keyId string) pgp {
	doc, _ := goquery.NewDocument(keyserver + keyId)
	keydata := &pgp{}
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
	return *keydata
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
func pypiClient(user string) []interface{} {
	client, _ := xmlrpc.NewClient(pypi, nil)
	defer client.Close()
	// PyPI user_packages()
	var result [][]string
	client.Call("user_packages", user, &result)
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
		client.Call("release_data", []interface{}{v[1], ver[0]}, &meta)
		pkgs[i] = meta
	}
	return pkgs
}

func mergeJson(a Account) []byte {
	js := simplejson.New()
	js.Set("deb", debPackages(a.DebianEmail))
	js.Set("udd", restClient(udd+"?email1="+a.DebianEmail+"&format=json").MustArray())
	js.Set("github", restClient(github+"/"+a.GithubUser+"/events"))
	js.Set("pypi", pypiClient(a.PypiUser))
	js.Set("pgp", pgpData(a.KeyId))
	data, _ := js.EncodePretty()
	return data
}

func main() {
	c := flag.String("c", "glaneuses.conf", "Configuration file")
	o := flag.String("o", "glaneuses.json", "Output file")
	flag.Parse()

	acct := readConfig(*c)
	err := ioutil.WriteFile(*o, mergeJson(acct), 0644)
	if err != nil {
		panic(err)
	}
}
