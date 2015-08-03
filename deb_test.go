package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const debbody = `<html>
<head>
<title>Debian Maintainer Dashboard</title>
</head>
<body>
<div id="content">
<h2 id='versions'>Versions</h2>
<table class="tablesorter table table-bordered">
<thead>
<tr>
    <th>source</th>
    <th>squeeze</th>
    <th>wheezy</th>
    <th>jessie</th>
    <th>sid</th>
    <th>experimental</th>
    <th>vcs</th>
    <th>upstream</th>
</tr>
</thead>
<tbody>
    <tr><td><a href="https://tracker.debian.org/example0" title="Maintained by Alice">example0</a></td>
       <td>
       &nbsp;
      </td>
       <td>
      0.3.3-1 &nbsp;
      </td>
        <td colspan="2" class="center">
      0.5.1-1 &nbsp;
      </td>
        <td colspan="2" class="center">
       &nbsp;
      </td>
    <td>
          <span class="prio_high" title="Newer upstream version available">0.5.3</a>
      &nbsp;
    </td>
    </tr>
    <tr><td><a href="https://tracker.debian.org/example1" title="Maintained by Alice">example1</a></td>
        <td colspan="2" class="center">
       &nbsp;
      </td>
        <td colspan="2" class="center">
      0.8-1 &nbsp;
      </td>
        <td colspan="2" class="center">
       &nbsp;
      </td>
    <td>
          <span class="prio_high" title="Newer upstream version available">0.8.2</a>
      &nbsp;
    </td>
    </tr>
</tbody>
</table>
</div>
</body>
</html>
`

func TestDebPackages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, debbody)
	}))
	defer ts.Close()

	var srv = service{"debian", "guest@example.org", ""}
	srv.uri = ts.URL
	d, err := srv.debPackages()
	if err != nil {
		t.Fatalf("%v", err)
	}
	for i, v := range d {
		if v.(deb).Source != fmt.Sprintf("example%d", i) {
			t.Fatalf("%v, want: example%d\n", v.(deb).Source, i)
		}
		if v.(deb).URL != fmt.Sprintf("https://tracker.debian.org/example%d", i) {
			t.Fatalf("%s, want: https://tracker.debian.org/example%d\n", v.(deb).URL, i)
		}
	}
}
