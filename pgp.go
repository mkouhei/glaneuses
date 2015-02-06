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

func (a *account) pgpData(keyserverUrl string) (pgp, error) {
	keydata := &pgp{}
	doc, err := goquery.NewDocument(keyserverUrl + a.KeyId)
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
