package main

import (
	"fmt"
	"net/http"
)

type karmas map[string]int

func main() {
	k := karmas(make(map[string]int))
	http.Handle("/", k)

	http.ListenAndServe(":8080", nil)
}

func (k karmas) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	w.Write([]byte("fart"))
}
