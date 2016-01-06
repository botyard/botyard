package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/botyard/botyard/command"
)

func Test_HTTPGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	c, err := New("GET", ts.URL, "hello <user>")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	req, ok := c.Match("hello ")
	if have, want := ok, true; have != want {
		t.Errorf("Match have: %v want: %v", have, want)
	}

	args := req.([]*command.Argument)
	if have, want := len(args), 0; have != want {
		t.Errorf("Match have: %v want: %v", have, want)
	}

	ctx := context.Background()

	resp, err := c.Endpoint()(ctx, req)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	if resp == nil {
		t.Errorf("err: resp")
	}

	t.Logf("%v", resp)

}
