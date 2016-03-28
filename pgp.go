package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type pgp struct {
	Payload       string
	PublicKeyPath string
	VindexPath    string
}

func (a *account) ignoreUID(s string) bool {
	for _, u := range strings.Split(a.IgnoreUids, ",") {
		if strings.Contains(s, u) {
			return true
		}

	}
	return false
}

func (a *account) pgpData(keyserverURL string) (pgp, error) {
	keydata := &pgp{}
	doc, err := goquery.NewDocument(keyserverURL + a.KeyID)
	if err != nil {
		return *keydata, err
	}
	doc.Find("pre+hr+pre").Each(func(i int, s *goquery.Selection) {
		var p string
		if a.IgnoreUids != "" {
			for _, l := range strings.Split(s.Text(), "\n") {
				if !a.ignoreUID(l) {
					p += l + "\n"
				}
			}
		} else {
			p = s.Text()
		}
		keydata.Payload = strings.Replace(
			strings.Replace(p, "@", " at ", -1),
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
