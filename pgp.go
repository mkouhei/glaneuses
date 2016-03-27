package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type pgp struct {
	Payload       string
	PublicKeyPath string
	VindexPath    string
	ignoreUIDS    []string
}

func ignoreUID(ignoreUIDS []string, s string) bool {
	for _, u := range ignoreUIDS {
		if strings.Contains(s, u) {
			return true
		}
	}
	return false
}

func (srv *service) pgpData(ignoreUIDS []string) (pgp, error) {
	keydata := &pgp{}

	doc, err := goquery.NewDocument(srv.uri)
	if err != nil {
		return *keydata, err
	}
	doc.Find("pre+hr+pre").Each(func(i int, s *goquery.Selection) {
		var p string
		if len(ignoreUIDS) > 0 {
			for _, l := range strings.Split(s.Text(), "\n") {
				if !ignoreUID(ignoreUIDS, l) {
					p += l + "\n"
				}
			}
		} else {
			p = s.Text()
		}
		keydata.Payload = strings.Replace(
			strings.Replace(p, "@", " at ", -1),
			".", " dot ", -1)
		keydata.Payload = strings.Replace(keydata.Payload, " \n\n", "\n", -1)
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
