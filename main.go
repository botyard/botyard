package main

import (
	"github.com/botyard/botyard/lib"
	"github.com/botyard/botyard/lib/message"
	"github.com/botyard/botyard/lib/sync"

	"github.com/codegangsta/cli"
	"golang.org/x/net/context"

	"fmt"
	"io/ioutil"
	"log"
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
	log.Println("Start.")

	yamlCfgFilePath := c.String("config")
	if yamlCfgFilePath == "" {
		log.Fatal("The config file not found")
		return
	}

	yamlContent, err := ioutil.ReadFile(yamlCfgFilePath)
	if err != nil {
		log.Fatalf("Read config file. err:%v", err)
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
		log.Fatalf("Loader err:%v", err)
		return
	}

	if err := loader.OpenGateways(msgChannel); err != nil {
		log.Fatalf("Open gateway err:%v", err)
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

	log.Println("End.")
}
