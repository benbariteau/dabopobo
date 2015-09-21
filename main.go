package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/naoina/toml"

	"github.com/firba1/dabopobo/lib"
)

var configPath = flag.String("config", "", "Configuration file for dabopobo")

const defaultConfigName = ".dabopobo.toml"

func main() {
	flag.Parse()
	conf, err := parseConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading configuration:", err)
		os.Exit(1)
	}

	err = lib.Serve(
		conf.Redis.Address,
		conf.slackTokens(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type config struct {
	Slack []struct {
		Name  string
		Token string
	}
	Redis struct {
		Address string
	}
}

func (c config) slackTokens() (tokens []string) {
	for _, slack := range c.Slack {
		tokens = append(tokens, slack.Token)
	}
	return
}

func parseConfig(configPath string) (conf config, err error) {
	if configPath == "" {
		usr, err := user.Current()
		if err != nil {
			return conf, err
		}
		configPath = filepath.Join(usr.HomeDir, defaultConfigName)
	}

	f, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = toml.Unmarshal(buf, &conf)
	if err != nil {
		return
	}
	return conf, nil
}
