package lib

import (
	"fmt"
	"regexp"
	"time"

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
	fmt.Println("rtm start")
	conn, err := rtm.Dial(token)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Println("rtm connected")

	events := eventsChannel(conn)

	for event := range events {
		fmt.Printf("handling %v event\n", event.Type())
		switch event.(type) {
		case rtm.Message:
			e := event.(rtm.Message)
			fmt.Println("handling message", e)
			handleMessage(e, conn, s)
		}
	}
	return nil
}

func handleMessage(message rtm.Message, conn *rtm.Conn, s serverConfig) {
	for _, command := range s.commands {
		r := regexp.MustCompile(command.regex)
		matches := r.FindAllStringSubmatch(conn.UnescapeMessage(message.Text()), -1)
		if matches == nil {
			continue
		}
		text, err := command.handler(s, matches, conn.UserInfo(message.User()).Name)
		if err != nil {
			efmt.Eprintln(err)
		}

		if text != "" {
			conn.SendMessage(text, message.Channel())
		}

		break
	}
}

func eventsChannel(conn *rtm.Conn) <-chan rtm.Event {
	events := make(chan rtm.Event)
	go func() {
		for {
			event := conn.NextEvent()
			fmt.Println(event.Type(), event)
			events <- event
		}
	}()
	return events
}
