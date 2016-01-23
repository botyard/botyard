package http

import (
	"github.com/botyard/botyard/command"
	"github.com/botyard/botyard/command/matcher"
	"github.com/botyard/botyard/command/parse"
	"github.com/botyard/botyard/message"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HttpCommand struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Cmd    string `json:"cmd"`
	items  []*parse.Item
	client *http.Client
}

func New(method, url, cmd string) (*HttpCommand, error) {
	c := &HttpCommand{
		Method: method,
		Url:    url,
		Cmd:    cmd,
		client: http.DefaultClient,
	}

	items, err := parse.Parse(c.Cmd)
	if err != nil {
		return nil, err
	}
	c.items = items
	return c, nil
}

func (c *HttpCommand) Match(in string) (req interface{}, ok bool) {
	m := matcher.New(in, c.items)
	args, ok := m.MatchAndReturnArguments()
	return args, ok
}

func (c *HttpCommand) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		u, err := url.Parse(c.Url)
		if err != nil {
			return nil, fmt.Errorf("URL: %v", err)
		}

		if c.Method == "GET" {
			q := u.Query()
			args, ok := req.([]*command.Argument)
			if !ok {
				return nil, fmt.Errorf("TypeError: The request is not argument")
			}

			if len(q) <= 0 {
				for _, arg := range args {
					q.Set(arg.Name, arg.Value)
				}
			} else {
				//Reverse query values
				//Filter query values e.g. message=<messageid>&user=<userid>
				kv := make(map[string]string)
				for key, val := range q {
					if strings.HasPrefix(val[0], "<") && strings.HasSuffix(val[0], ">") {
						kv[val[0][1:len(val[0])-1]] = key
					}
				}
				for _, arg := range args {
					param, ok := kv[arg.Name]
					if !ok {
						continue
					}
					q.Set(param, arg.Value)
				}
			}
			u.RawQuery = q.Encode()
		} else if c.Method == "POST" {
			//TODO: How to represent POST form
			//For now, There is same scheme for the command argument in the POST method

		}

		request, err := http.NewRequest(c.Method, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("NewRequest: %v", err)
		}

		response, err := ctxhttp.Do(ctx, c.client, request)
		if err != nil {
			return nil, fmt.Errorf("Do: %v", err)
		}
		//TODO:
		//defer func() { _ = response.Body.Close() }()

		return response, nil
	}
}

func (c *HttpCommand) Response(res interface{}) (*message.Message, error) {
	response, ok := res.(*http.Response)
	if !ok {
		return nil, fmt.Errorf("res is not http.Response")
	}
	defer func() { _ = response.Body.Close() }()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	msg := &message.Message{}
	msg.Body = string(body) //TODO
	log.Println(msg.Body)

	return msg, nil
}
