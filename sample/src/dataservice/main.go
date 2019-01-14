package main

import (
	"fmt"
	"os"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	id := r.URL.Path[len("/data/"):]
	fmt.Fprintf(w, "hello, this is data %s-%s", id, hostname)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9090", nil)
}
