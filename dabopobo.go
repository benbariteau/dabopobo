package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuyu/goredis"
)

type serverConfig struct {
	redis *goredis.Redis
}

type karmaSet struct {
	plusplus   int
	minusminus int
	plusminus  int
}

func (k karmaSet) value() int {
	return k.plusplus - k.minusminus
}

func (k karmaSet) String() string {
	return fmt.Sprintf("(%v++,%v--,%v+-)", k.plusplus, k.minusminus, k.plusminus)
}

var indentifierRegex = regexp.MustCompile("([^ ]+)(\\+\\+|--|\\+-|-\\+)")
var getkarma = regexp.MustCompile("^!karma +([^ ]+)")

var redisPort = flag.Int("redisport", 6379, "redis port")
var port = flag.Int("port", 8080, "port")

func main() {
	flag.Parse()
	err := serve()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve() error {
	redis, err := goredis.Dial(&goredis.DialConfig{Address: fmt.Sprintf("127.0.0.1:%v", *redisPort)})
	if err != nil {
		return err
	}
	s := serverConfig{redis}

	http.Handle("/", s)

	return http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
}

func (s serverConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	indentifierMatches := indentifierRegex.FindAllStringSubmatch(text, -1)
	karma := getkarma.FindStringSubmatch(text)
	username := r.Form.Get("user_name")
	if karma != nil {
		name := karma[1]
		fmt.Println("asking for", name)
		res := make(map[string]string)
		karmaset := s.getKarmaSet(name)
		fmt.Println(karmaset)
		res["text"] = fmt.Sprintf("%v's karma is %v %v", name, karmaset.value(), karmaset)
		res["parse"] = "full"
		res["username"] = "dabopobo"
		resp, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resp))
		w.WriteHeader(200)
		w.Write(resp)
	} else if indentifierMatches != nil && username != "slackbot" {
		for _, match := range indentifierMatches {
			key := match[1]
			op := match[2]
			if key != "" && key != username {
				suffix := canonicalizeSuffix(op)
				_, err := s.redis.Incr(strings.ToLower(key) + suffix)
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func canonicalizeSuffix(suffix string) string {
	switch suffix {
	case "--", "++", "+-":
		return suffix
	case "-+":
		return "+-"
	default:
		return suffix
	}
}

func (s serverConfig) getKarmaSet(name string) (k karmaSet) {
	name = strings.ToLower(name)
	k.plusplus = getRedisInt(s.redis, name+"++", 0)
	k.minusminus = getRedisInt(s.redis, name+"--", 0)
	k.plusminus = getRedisInt(s.redis, name+"+-", 0)
	return
}

func getRedisInt(r *goredis.Redis, key string, def int) int {
	val, err := r.Get(key)
	if err != nil {
		return def
	}

	value, err := strconv.Atoi(string(val))
	if err != nil {
		return def
	}

	return value
}
