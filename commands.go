package main

import (
	"strings"
)

type model interface {
	incr(string) error
}

type cmd struct {
	regex string
}

func mutateKarma(m model, mutations [][]string, username string) (err error) {
	if username == "slackbot" {
		return nil
	}
	for _, mutation := range mutations {
		identifier := mutation[1]
		op := mutation[2]
		if identifier != "" && identifier != username { //users may not mutate themselves
			suffix := canonicalizeSuffix(op)
			err = m.incr(strings.ToLower(identifier) + suffix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
