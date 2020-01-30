package main

import (
	"fmt"
	"net/http"
)

const (
	addr = ":8080"
)

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("SubmitHandler:%v\n", r.Form)
}

func main() {
	http.HandleFunc("/submit", SubmitHandler)
	fmt.Printf("start serv on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
