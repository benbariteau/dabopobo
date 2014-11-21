package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/xuyu/goredis"
)

var redisPort = flag.Int("redisport", 6379, "redis port")
var port = flag.Int("port", 8080, "port")

func main() {
	flag.Parse()
	err := serve()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve() error {
	redis, err := goredis.Dial(&goredis.DialConfig{Address: fmt.Sprintf("127.0.0.1:%v", *redisPort)})
	if err != nil {
		return err
	}
	s := serverConfig{
		redis,
		[]cmd{
			cmd{"^!karma +([^ ]+)", getKarma},
			cmd{"([^ ]+)(\\+\\+|--|\\+-|-\\+)", mutateKarma},
		},
	}

	http.Handle("/", s)

	return http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
}

func (s serverConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	username := r.Form.Get("user_name")
	for _, command := range s.commands {
		r := regexp.MustCompile(command.regex)
		matches := r.FindAllStringSubmatch(text, -1)
		if matches == nil {
			continue
		}
		response, err := command.handler(s, matches, username)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if len(response) > 0 {
			w.Write(response)
		}
	}
}

func canonicalizeSuffix(suffix string) string {
	switch suffix {
	case "--", "++", "+-":
		return suffix
	case "-+":
		return "+-"
	default:
		return suffix
	}
}
