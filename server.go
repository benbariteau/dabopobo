package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/xuyu/goredis"
)

func serve(port uint16, redisPort uint16) error {
	redis, err := goredis.Dial(&goredis.DialConfig{Address: fmt.Sprintf("127.0.0.1:%v", redisPort)})
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

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func (s serverConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	username := r.Form.Get("user_name")

	// process commands in order so that there is an order of precedence
	// this is to prevent mutation on query, for example
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

		w.Write(response)
	}
}
