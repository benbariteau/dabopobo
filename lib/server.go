package lib

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/firba1/slack/rtm"
	_ "github.com/mattn/go-sqlite3"
)

func Serve(redisAddr string, slackTokens []string) error {
	commands := []cmd{
		rawGetKarmaCmd,
		getKarmaCmd,
		helpCmd,        // must be after getKarma since its regex matches anything that getKarma's does
		mutateKarmaCmd, // must be after getKarma because getKarma could be given an identifier that matches
		singleMessageKarmaMutate,
		mentionedCmd,
		leaderboardCmd,
	}
	db, err := sql.Open("sqlite3", "./dabopobo.db")
	if err != nil {
		return err
	}
	s := sqliteBackend{
		db: db,
	}

	for _, slackToken := range slackTokens {
		go rtmHandle(slackToken, commands, s)
	}
	for {
		time.Sleep(5 * time.Minute)
	}
	return nil
}

func rtmHandle(token string, commands []cmd, m model) error {
	for {
		fmt.Println("rtm start")
		conn, err := rtm.Dial(token)
		if err != nil {
			continue
		}
		defer conn.Close()

		fmt.Println("rtm connected")

		rtmHelper(conn, commands, m)
		fmt.Println("attempting rtm reconnect")
	}
	return nil
}

func rtmHelper(conn *rtm.Conn, commands []cmd, m model) {
	events := eventsChannel(conn)
	ticker := time.Tick(10 * time.Minute)

	eventHandledRecently := true
	for {
		select {
		case event := <-events:
			handleEvent(event, conn, commands, m)
			eventHandledRecently = true
		case <-ticker:
			if eventHandledRecently {
				eventHandledRecently = false
			} else {
				return
			}
		}
	}
}

func handleEvent(event rtm.Event, conn *rtm.Conn, commands []cmd, m model) {
	fmt.Printf("handling %v event\n", event.Type())
	switch event.(type) {
	case rtm.Message:
		e := event.(rtm.Message)
		fmt.Println("handling message", e)
		handleMessage(e, conn, commands, m)
	}
}

func postprocessUser(userString string) string {
	// string leading "@"
	userString = userString[1:]

	// wrap in parens so that it gets treated as a single unit for karma purposes, even if it has spaces
	return "(" + userString + ")"
}

func postprocessLink(unescapedString string, escapeType rtm.EscapeType) string {
	if escapeType != rtm.UserEscape {
		return unescapedString
	}

	return postprocessUser(unescapedString)
}

func handleMessage(message rtm.Message, conn *rtm.Conn, commands []cmd, m model) {
	for _, command := range commands {
		r := regexp.MustCompile(command.regex)
		unescapedText := conn.UnescapeMessagePostprocess(message.Text(), postprocessLink)
		matches := r.FindAllStringSubmatch(unescapedText, -1)
		if matches == nil {
			continue
		}
		fmt.Printf("matched with %v\n", command)
		text, err := command.handler(m, matches, conn.UserInfo(message.User()).Profile.DisplayName, message.Channel())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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
