package main

import (
	"github.com/botyard/botyard/lib"
	"github.com/botyard/botyard/lib/log"
	"github.com/botyard/botyard/lib/message"
	"github.com/botyard/botyard/lib/sync"

	"github.com/codegangsta/cli"
	"golang.org/x/net/context"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "botyard"
	app.Usage = "botyard - make bots,together!"
	app.Version = fmt.Sprintf("%s commit %s build %s", version, gitCommit, buildTime)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,c",
			Value:  "config.yaml",
			EnvVar: "BOTYARD_CONFIG",
		},
		cli.IntFlag{
			Name:   "num_dispatcher,n",
			Value:  1,
			EnvVar: "BOTYARD_NUM_D",
		},
	}
	app.Action = DaemonAction
	app.Run(os.Args)
}

func DaemonAction(c *cli.Context) {
	logger := log.Logger
	logger.Log("start", 1)

	yamlCfgFilePath := c.String("config")
	if yamlCfgFilePath == "" {
		logger.Log("err", "config file not found")
		return
	}

	yamlContent, err := ioutil.ReadFile(yamlCfgFilePath)
	if err != nil {
		logger.Log("err", err, "config", "file")
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer func() { cancelFunc() }()

	var msgChannel chan *message.Message
	{
		msgChannel = make(chan *message.Message)
	}

	loader, err := lib.NewLoader(yamlContent)
	if err != nil {
		logger.Log("err", err, "loader", 0)
		return
	}

	if err := loader.OpenGateways(msgChannel); err != nil {
		logger.Log("err", err, "gateway", 0)
		return
	}

	var dispatchers []*lib.Dispatcher
	{
		for i := 0; i < c.Int("num_dispatcher"); i++ {
			dispatcher := lib.NewDispatcher(ctx, msgChannel, loader)
			dispatchers = append(dispatchers, dispatcher)
			sync.WaitGroup.Add(1)
		}
	}

	go func() {
		http.HandleFunc("/hello", lib.HelloHandler)
		http.ListenAndServe(":9001", nil)
	}()
	sync.WaitGroup.Wait()

	logger.Log("end", 1)
}
