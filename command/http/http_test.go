package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/botyard/botyard/command"
)

func Test_HTTPGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("u")
		fmt.Fprintln(w, fmt.Sprintf("Hello, %v", user))
	}))
	defer ts.Close()

	c, err := New("GET", ts.URL+"?u=<user>", "hello <user>")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	req, ok := c.Match("hello client")
	if have, want := ok, true; have != want {
		t.Errorf("Match have: %v want: %v", have, want)
	}

	args := req.([]*command.Argument)
	if have, want := len(args), 1; have != want {
		t.Errorf("Match args have: %v want: %v", have, want)
	}
	t.Logf("args:%v", args)

	ctx := context.Background()

	resp, err := c.Endpoint()(ctx, req)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	if resp == nil {
		t.Errorf("err: resp")
	}

	response := resp.(*http.Response)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	//t.Logf("%v", resp)
	t.Logf("body:%v", string(body))
	if have, want := string(body), "Hello, client\n"; have != want {
		t.Errorf("Response have: %v want: %v", have, want)
	}

}
