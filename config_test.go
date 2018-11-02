package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	conf := &config{}
	conf.loadConfig("examples/glaneuses.conf")
	for _, srv := range conf.services {
		if srv.name == "debian" && srv.account != "guest@example.org" {
			t.Fatal("parse error [debian]email")
		}
		if srv.name == "pypi" && srv.account != "guest_pypi" {
			t.Fatal("parse error [pypi]username")
		}
		if srv.name == "rubygems" && srv.account != "guest_gems" {
			t.Fatal("parse error [rubygems]username")
		}
		if srv.name == "github" && srv.account != "guest_github" {
			t.Fatal("parse error [github]username")
		}
		if srv.name == "bitbucket" && srv.account != "guest_bitbucket" {
			t.Fatal("parse error [bitbucket]username")
		}
	}
}
