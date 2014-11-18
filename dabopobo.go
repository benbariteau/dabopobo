package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type karmas struct {
	karmaMap *map[string]int
}

var regex = regexp.MustCompile("([^ ]+)\\+\\+")
var getkarma = regexp.MustCompile("^!karma +([^ ]+)")

func main() {
	k := newKarmas()
	http.Handle("/", k)

	http.ListenAndServe(":8080", nil)
}

func newKarmas() karmas {
	m := make(map[string]int)
	return karmas{&m}
}

func (k karmas) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	matches := regex.FindAllStringSubmatch(text, -1)
	karma := getkarma.FindStringSubmatch(text)
	if matches != nil {
		for _, match := range matches {
			key := match[1]
			if key != "" {
				(*k.karmaMap)[key]++
			}
		}
		fmt.Println(*k.karmaMap)
	} else if karma != nil {
		name := karma[1]
		fmt.Println("asking for", name)
		res := make(map[string]string)
		res["text"] = fmt.Sprintf("%v: %v", r.Form.Get("user_name"), (*k.karmaMap)[name])
		res["parse"] = "full"
		resp, _ := json.Marshal(res)
		fmt.Println(string(resp))
		w.WriteHeader(200)
		w.Write(resp)
	}
}
