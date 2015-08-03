package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type deb struct {
	Source string
	URL    string
}

func (srv *service) debPackages() ([]interface{}, error) {
	doc, err := goquery.NewDocument(srv.uri)
	if err != nil {
		return nil, err
	}
	cnt := doc.Find("h2#versions+table a").Length()
	debs := make([]interface{}, cnt)
	doc.Find("h2#versions+table a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists && !strings.Contains(url, "qa.debian.org") {
			debs[i] = deb{s.Text(), url}
		}
	})
	return debs, nil
}
