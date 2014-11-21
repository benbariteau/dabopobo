package main

import (
	"github.com/xuyu/goredis"
	"strconv"
)

type serverConfig struct {
	redis    *goredis.Redis
	commands []cmd
}

func (s serverConfig) incr(key string) error {
	_, err := s.redis.Incr(key)
	return err
}

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
