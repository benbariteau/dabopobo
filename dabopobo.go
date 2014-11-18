package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type karmas struct {
	karmaMap *map[string]karmaSet
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

var regex = regexp.MustCompile("([^ ]+)(\\+\\+|--|\\+-|-\\+)")
var getkarma = regexp.MustCompile("^!karma +([^ ]+)")

func main() {
	k := newKarmas()
	http.Handle("/", k)

	http.ListenAndServe(":8080", nil)
}

func newKarmas() karmas {
	m := make(map[string]karmaSet)
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
			op := match[2]
			set := (*k.karmaMap)[key]
			if key != "" {
				switch op {
				case "--":
					set.minusminus++
				case "++":
					set.plusplus++
				case "-+", "+-":
					set.plusminus++
				}
				(*k.karmaMap)[key] = set
			}
		}
		fmt.Println(*k.karmaMap)
	} else if karma != nil {
		name := karma[1]
		fmt.Println("asking for", name)
		res := make(map[string]string)
		karmaset := (*k.karmaMap)[name]
		res["text"] = fmt.Sprintf("%v: %v %v", r.Form.Get("user_name"), karmaset.value(), karmaset)
		res["parse"] = "full"
		res["username"] = "dabopobo"
		resp, _ := json.Marshal(res)
		fmt.Println(string(resp))
		w.WriteHeader(200)
		w.Write(resp)
	}
}
