package lib

import (
	"fmt"
	"strings"
)

type model interface {
	incr(key string) error                                        //increment the given key, setting it to zero if it doesn't exist
	getInt(key string) int                                        //get a key as an int, defaulting to 0 if it doesn't exist
	addChannelKarma(mutation karmaMutation, channel string) error // apply given karma mutation to channel
}

// a dabopobo command, consisting of a regex to match against and a commandHandler to run if it matches
type cmd struct {
	regex   string
	handler commandHandler
}

/*
A commandHandler is a function that handles a dabopobo command and returns a slice of bytes to respond with and possibly an error.

m is the data model for this process.
submatches is the output of FindAllStringSubmatch() when running the message text on the corresponding cmd's regex.
username is the username of the user who sent the message.
response is the response to be sent back to slack.
err is non-nil if an error is produced. response should be empty in this case.
*/
type commandHandler func(m model, submatches [][]string, username string, channel string) (response string, err error)

var mutateKarmaCmd = cmd{"(\\(.+\\)|@?[^ ]+?)(\\+\\++|--+|\\+-|-\\+)", mutateKarma}
var singleMessageKarmaMutate = cmd{`^(\(.+\)|@?[^ ]+?)\s+(\+\++|--+|\+-|-\+)$`, mutateKarma}

//handles identifier++
func mutateKarma(m model, mutations [][]string, username string, channel string) (s string, err error) {
	karmaMutations := filterMutations(mutations, username)
	for _, mutation := range karmaMutations {
		key := mutation.key()
		err = m.addChannelKarma(mutation, channel)
		if err != nil {
			err = nil
			return
		}
		err = m.incr(key)
		if err != nil {
			err = nil
			return
		}
		fmt.Println(key)
	}
	return
}

var rawGetKarmaCmd = cmd{`^!karma\s+\((.*)\)`, getKarma}
var getKarmaCmd = cmd{"^!karma +([^ ].*)", getKarma}

//handles !karma identifier
func getKarma(m model, identifier [][]string, username string, channel string) (text string, err error) {
	name := identifier[0][1] //since the regex has a beginning of string hook, there should only be one match, so we only care about index 0.
	karmaset := getKarmaSet(m, name)
	text = fmt.Sprintf("%v's karma is %v %v", name, karmaset.value(), karmaset)
	fmt.Println(text)
	return
}

var helpCmd = cmd{"^!karma(|help)$", help}

func help(m model, s [][]string, u string, channel string) (string, error) {
	fmt.Println("help message")
	return strings.Join(
		[]string{
			"thing++\t(with at least 2 pluses) gives positive karma",
			"thing--\t(with at least 2 minuses) gives negative karma",
			"thing+-\t(either order) gives neutral karma",
			"(thing with spaces)++\t(with any of the above) gives karma to a thing with spaces in it",
			"!karma thing\tdisplays things karma. It can be anything, including with spaces in the middle",
			"!karma (any string )\tdisplays karma for \"any string \", may include trailng spaces",
			"!karma or !karmahelp\tdisplays this help",
		},
		"\n",
	), nil
}

var mentionedCmd = cmd{"dabopobo", mentioned}

func mentioned(m model, s [][]string, u string, channel string) (string, error) {
	fmt.Println("don't touch me")
	return "don't touch me", nil
}

// getKarmaSet loads the karma for a given key
func getKarmaSet(m model, name string) (k karmaSet) {
	name = strings.ToLower(name)
	k.plusplus = m.getInt(name + "++")
	k.minusminus = m.getInt(name + "--")
	k.plusminus = m.getInt(name + "+-")
	return
}
