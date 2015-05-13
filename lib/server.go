package lib

import (
	"regexp"
	"time"

	"github.com/firba1/slack"
	"github.com/firba1/slack/rtm"
	"github.com/firba1/util/efmt"
	"github.com/xuyu/goredis"
)

func Serve(redisAddr string, slackTokens []string) error {
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

	for _, slackToken := range slackTokens {
		go rtmHandle(slackToken, s)
	}
	for {
		time.Sleep(5 * time.Minute)
	}
	return nil
}

func rtmHandle(token string, s serverConfig) error {
	conn, err := rtm.Dial(token)
	if err != nil {
		return err
	}
	defer conn.Close()

	messages := conn.MessageChan()

	slackAPI := slack.NewAPI(token)

	for message := range messages {
		userInfo, err := slackAPI.UsersInfo(message.User)
		if err != nil {
			efmt.Eprintln(err)
			continue
		}

		for _, command := range s.commands {
			r := regexp.MustCompile(command.regex)
			matches := r.FindAllStringSubmatch(message.Text, -1)
			if matches == nil {
				continue
			}
			text, err := command.handler(s, matches, userInfo.Name)
			if err != nil {
				efmt.Eprintln(err)
			}

			if text != "" {
				conn.SendMessage(text, message.Channel)
			}

			break
		}
	}
	return nil
}
