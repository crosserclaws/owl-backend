package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Cepave/query/conf"
	"github.com/Cepave/query/database"
	"github.com/Cepave/query/g"
	ginHttp "github.com/Cepave/query/gin_http"
	"github.com/Cepave/query/graph"
	"github.com/Cepave/query/grpc"
	"github.com/Cepave/query/http"
	"github.com/Cepave/query/proc"
)

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
	version := flag.Bool("v", false, "show version")
	versionGit := flag.Bool("vg", false, "show version and git commit log")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if *versionGit {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// config
	g.ParseConfig(*cfg)
	// proc
	proc.Start()

	// graph
	graph.Start()

	grpcMsg := make(chan string)
	if gconf.Grpc.Enabled {
		// grpc
		go grpc.Start(grpcMsg)
	}


	//lambdaSetup
	database.Init()
	conf.ReadConf("./conf/lambdaSetup.json")
	go ginHttp.StartWeb()
	select {}
	httpMsg := make(chan string)

	if gconf.Http.Enabled {
		// http
		go http.Start(httpMsg)
	}
}
