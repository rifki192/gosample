package main

import (
	"flag"
	"github.com/google/gops/agent"
	"log"
	"net/http"

	"github.com/tokopedia/gosample/hello"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	// gops helps us get stack trace if something wrong/slow in production
	if err := agent.Listen(agent.Options{ShutdownCleanup: true}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule()

	http.HandleFunc("/hello", hwm.SayHelloWorld)
	go logging.StatsLog()

	log.Fatal(grace.Serve(":9000", nil))
}
