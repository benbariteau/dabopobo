package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/firba1/util/efmt"
	"github.com/naoina/toml"

	"github.com/firba1/dabopobo/lib"
)

var configPath = flag.String("config", "", "Configuration file for dabopobo")

const defaultConfigName = ".dabopobo.toml"

func main() {
	flag.Parse()
	conf, err := parseConfig(*configPath)
	if err != nil {
		efmt.Fatalln("Error loading configuration:", err)
	}

	err = lib.Serve(
		conf.Redis.Address,
		conf.Slack.Token,
	)
	if err != nil {
		efmt.Fatalln(err)
	}
}

type config struct {
	Slack struct {
		Token string
	}
	Redis struct {
		Address string
	}
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
