package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/xuyu/goredis"
)

func serve(port uint16, redisNetwork, redisAddr string) error {
	redis, err := goredis.Dial(&goredis.DialConfig{Network: redisNetwork, Address: redisAddr})
	if err != nil {
		return err
	}
	s := serverConfig{
		redis,
		[]cmd{
			getKarmaCmd,
			helpCmd,        // must be after getKarma since its regex matches anything that getKarma's does
			mutateKarmaCmd, // must be after getKarma because getKarma could be given an identifier that matches
			mentionedCmd,
		},
	}

	http.Handle("/", s)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func (s serverConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	username := r.Form.Get("user_name")
	if username == "slackbot" { // ignore messages from bot(s)
		return
	}

	// process commands in order so that there is an order of precedence
	// this is to prevent mutation on query, for example
	for _, command := range s.commands {
		r := regexp.MustCompile(command.regex)
		matches := r.FindAllStringSubmatch(text, -1)
		if matches == nil {
			continue
		}
		text, err := command.handler(s, matches, username)
		var response []byte
		if text != "" {
			res := map[string]string{
				"text":     text,
				"parse":    "full",     // allows the user to be pinged
				"username": "dabopobo", // so IRC users don't some weird thing because slack
			}
			response, err = json.Marshal(res)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		w.Write(response)
		return
	}
}
