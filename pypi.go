package main

import (
	"github.com/kolo/xmlrpc"
)

type dl struct {
	LastDay   int `xmlrpc:"last_day"`
	LastMonth int `xmlrpc:"last_month"`
	LastWeek  int `xmlrpc:"last_week"`
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
func (a *account) pypiClient() ([]interface{}, error) {
	client, err := xmlrpc.NewClient(pypi, nil)
	if err != nil {
		return nil, err
	}

	// PyPI user_packages()
	var result [][]string
	client.Call("user_packages", a.PypiUser, &result)
	pkgs := make([]interface{}, len(result))
	for i, v := range result {
		var ver []string
		err = client.Call("package_releases", v[1], &ver)
		if err != nil {
			return nil, err
		}
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
			client.Close()
			return nil, err
		}
		pkgs[i] = meta
	}

	client.Close()
	return pkgs, nil
}
