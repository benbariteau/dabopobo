package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type karmas map[string]int

var regex = regexp.MustCompile("[^ ]+\\+\\+")

func main() {
	k := karmas(make(map[string]int))
	http.Handle("/", k)

	http.ListenAndServe(":8080", nil)
}

func (k karmas) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.Form.Get("text")
	matches := regex.FindAllString(text, -1)
	if matches == nil {
		return
	}
	for _, match := range matches {
		k[match]++
	}
}
