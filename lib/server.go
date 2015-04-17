package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/firba1/slack/rtm"
	"github.com/firba1/util/efmt"
	"github.com/xuyu/goredis"
)

func Serve(port uint16, redisAddr string, slackToken string) error {
	redis, err := goredis.Dial(&goredis.DialConfig{Address: redisAddr})
	if err != nil {
		return err
	}
	s := serverConfig{
		redis,
		[]cmd{
			getKarmaCmd,
			helpCmd, // must be after getKarma since its regex matches anything that getKarma's does
			//mutateKarmaCmd, // must be after getKarma because getKarma could be given an identifier that matches
			mentionedCmd,
		},
	}

	go rtmHandle(slackToken, serverConfig{redis, []cmd{mutateKarmaCmd}})

	http.Handle("/", s)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func rtmHandle(token string, s serverConfig) {
	conn, err := rtm.Dial(token)
	if err != nil {
		efmt.Fatalln("Unable to connect to slack websocket:", err)
	}
	defer conn.Close()

	messages := conn.MessageChan()

	for message := range messages {
		for _, command := range s.commands {
			r := regexp.MustCompile(command.regex)
			matches := r.FindAllStringSubmatch(message.Text, -1)
			if matches == nil {
				continue
			}
			text, err := command.handler(s, matches, message.User)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			if text != "" {
				efmt.Eprintln("trying to reply: ", text)
			}

			return
		}
	}
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
