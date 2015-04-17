package lib

import (
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
			helpCmd,        // must be after getKarma since its regex matches anything that getKarma's does
			mutateKarmaCmd, // must be after getKarma because getKarma could be given an identifier that matches
			mentionedCmd,
		},
	}

	rtmHandle(slackToken, s)

	return nil
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
				efmt.Eprintln(err)
			}

			if text != "" {
				conn.SendMessage(text, message.Channel)
			}

			break
		}
	}
}
