package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const body = `
<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd" >
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>Search results for '0x00000000'</title>
<meta http-equiv="Content-Type" content="text/html;charset=utf-8" />
<style type="text/css">
/*<![CDATA[*/
 .uid { color: green; text-decoration: underline; }
 .warn { color: red; font-weight: bold; }
/*]]>*/
</style></head><body><h1>Search results for '0x00000000'</h1><pre>Type bits/keyID     Date       User ID
</pre><hr /><pre>
pub  4096R/<a href="/pks/lookup?op=get&amp;search=0x0000000000000000">00000000</a> 2014-08-12 <a href="/pks/lookup?op=vindex&amp;fingerprint=on&amp;search=0x0000000000000000">Alice &lt;alice@example.org&gt;</a>
	 Fingerprint=0000 0000 0000 0000 0000  0000 0000 0000 0000 0000 
</pre></body></html>

`

func TestPgpData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, body)
	}))
	defer ts.Close()

	srv := service{"pgp", "0x00000000", ""}
	srv.uri = ts.URL
	p, err := srv.pgpData([]string{"ignore@example.com"})
	if err != nil {
		t.Fatalf("%v", err)
	}
	payload := `pub  4096R/00000000 2014-08-12 Alice <alice at example dot org>
	 Fingerprint=0000 0000 0000 0000 0000  0000 0000 0000 0000 0000
`
	if p.Payload != payload {
		t.Fatalf("error: want: %s", payload)
	}
	pkeypath := "/pks/lookup?op=get&search=0x0000000000000000"
	if p.PublicKeyPath != pkeypath {
		t.Fatalf("error: want: %s", pkeypath)
	}
	vidxpath := "/pks/lookup?op=vindex&fingerprint=on&search=0x0000000000000000"
	if p.VindexPath != vidxpath {
		t.Fatalf("error: want: %s", vidxpath)
	}
}
