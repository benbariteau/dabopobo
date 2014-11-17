package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		w.Write([]byte("fart"))
	})

	http.ListenAndServe(":8080", nil)
}
