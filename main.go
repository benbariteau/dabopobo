package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/firba1/dabopobo/lib"
)

var redisAddr = flag.String("redisaddr", "127.0.0.1:6379", "redis backend port")
var slackAccessToken = flag.String("slackToken", "", "Accss token for corresponding slack bot")
var port = flag.Uint("port", 8080, "port")

func main() {
	flag.Parse()
	err := lib.Serve(uint16(*port), *redisAddr, *slackAccessToken)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
