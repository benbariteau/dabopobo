package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type model interface {
	incr(string) error
	getInt(string) int
}

type cmd struct {
	regex string
}

func mutateKarma(m model, mutations [][]string, username string) (b []byte, err error) {
	if username == "slackbot" {
		return
	}
	for _, mutation := range mutations {
		identifier := mutation[1]
		op := mutation[2]
		if identifier != "" && identifier != username { //users may not mutate themselves
			suffix := canonicalizeSuffix(op)
			err = m.incr(strings.ToLower(identifier) + suffix)
			if err != nil {
				err = nil
				return
			}
		}
	}
	return
}

func getKarma(m model, identifier [][]string, username string) (response []byte, err error) {
	name := identifier[0][1]
	fmt.Println("asking for", name)
	karmaset := getKarmaSet(m, name)
	res := map[string]string{
		"text":     fmt.Sprintf("%v's karma is %v %v", name, karmaset.value(), karmaset),
		"parse":    "full",
		"username": "dabopobo",
	}
	response, err = json.Marshal(res)
	return
}

func getKarmaSet(m model, name string) (k karmaSet) {
	name = strings.ToLower(name)
	k.plusplus = m.getInt(name + "++")
	k.minusminus = m.getInt(name + "--")
	k.plusminus = m.getInt(name + "+-")
	return
}
