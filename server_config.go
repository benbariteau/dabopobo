package main

import (
	"github.com/xuyu/goredis"
	"strconv"
)

// serverConfig implements the model interface, allowing for testablity
type serverConfig struct {
	redis    *goredis.Redis // redis is the redis backend for dabopobo
	commands []cmd          // commmands is a list of cmds in order of precedence, that are available
}

func (s serverConfig) incr(key string) error {
	_, err := s.redis.Incr(key)
	return err
}

// gets an int value from redis and returns zero if any errors occur (missing key, not an int, etc)
func (s serverConfig) getInt(key string) int {
	val, err := s.redis.Get(key)
	if err != nil {
		return 0
	}

	value, err := strconv.Atoi(string(val))
	if err != nil {
		return 0
	}

	return value
}
