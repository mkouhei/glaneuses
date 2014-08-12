package main

import (
	"github.com/PuerkitoBio/goquery"
)

type deb struct {
	Source string
	Url    string
}

func (a *Account) debPackages(uddUrl string) ([]interface{}, error) {
	doc, err := goquery.NewDocument(uddUrl + "?email1=" + a.DebianEmail)
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
