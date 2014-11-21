package main

import (
	"flag"
	"fmt"
	"os"
)

var redisPort = flag.Uint("redisport", 6379, "redis port")
var port = flag.Uint("port", 8080, "port")

func main() {
	flag.Parse()
	err := serve(uint16(*port), uint16(*redisPort))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
