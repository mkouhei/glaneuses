package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	a := &Account{}
	a.readConfig("examples/glaneuses.conf")
	if a.DebianEmail != "guest@example.org" {
		t.Fatal("parse error [debian]email")
	}
	if a.PypiUser != "guest" {
		t.Fatal("parse error [pypi]username")
	}
	if a.GemsUser != "guest" {
		t.Fatal("parse error [rubygems]username")
	}
	if a.GithubUser != "guest" {
		t.Fatal("parse error [github]username")
	}
	if a.BitbucketUser != "guest" {
		t.Fatal("parse error [bitbucket]username")
	}
	if a.KeyId != "0x00000000" {
		t.Fatal("parse error [pgp]keyid")
	}
}
