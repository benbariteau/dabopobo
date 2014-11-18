package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/xuyu/goredis"
)

type state struct {
	redis *goredis.Redis
}

type karmaSet struct {
	plusplus    int
	minusminus  int
	plusminus   int
	minusequals int
	plusequals  int
}

func (k karmaSet) value() int {
	return k.plusplus - k.minusminus
}

func (k karmaSet) String() string {
	return fmt.Sprintf("(%v++,%v--,%v+-,%v+=,%v-=)", k.plusplus, k.minusminus, k.plusminus, k.plusequals, k.minusequals)
}

var indentifierRegex = regexp.MustCompile("([^ ]+)(\\+\\+|--|\\+-|-\\+|-=|\\+=)")
var getkarma = regexp.MustCompile("^!karma +([^ ]+)")

func main() {
	redis, err := goredis.Dial(&goredis.DialConfig{Address: "127.0.0.1:6379"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	s := state{redis}
	http.Handle("/", s)

	http.ListenAndServe(":8080", nil)
}

func (s state) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	indentifierMatches := indentifierRegex.FindAllStringSubmatch(text, -1)
	karma := getkarma.FindStringSubmatch(text)
	username := r.Form.Get("user_name")
	if indentifierMatches != nil && username != "slackbot" {
		for _, match := range indentifierMatches {
			key := match[1]
			op := match[2]
			if key != "" && key != username {
				var err error
				switch op {
				case "--":
					_, err = s.redis.Incr(key + "--")
				case "++":
					_, err = s.redis.Incr(key + "++")
				case "-+", "+-":
					_, err = s.redis.Incr(key + "+-")
				case "+=":
					_, err = s.redis.Incr(key + "+=")
				case "-=":
					_, err = s.redis.Incr(key + "-=")
				}
				fmt.Fprintln(os.Stderr, err)
				if err != nil {
					panic(err)
				}
			}
		}
	} else if karma != nil {
		name := karma[1]
		fmt.Println("asking for", name)
		res := make(map[string]string)
		karmaset := s.getKarmaSet(name)
		res["text"] = fmt.Sprintf("%v's karma is %v %v", name, karmaset.value(), karmaset)
		res["parse"] = "full"
		res["username"] = "dabopobo"
		resp, _ := json.Marshal(res)
		fmt.Println(string(resp))
		w.WriteHeader(200)
		w.Write(resp)
	}
}

func (s state) getKarmaSet(name string) (k karmaSet) {
	k.plusplus = getRedisInt(s.redis, name+"++", 0)
	k.minusminus = getRedisInt(s.redis, name+"--", 0)
	k.plusminus = getRedisInt(s.redis, name+"+-", 0)
	k.plusequals = getRedisInt(s.redis, name+"+=", 0)
	k.minusequals = getRedisInt(s.redis, name+"-=", 0)
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
