package main

import (
	"fmt"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/api/cursor", func(w http.ResponseWriter, r *http.Request) { CursorWebsocketHandler(w, r, hub) })

	http.Handle("/", http.FileServer(http.Dir("./public")))

	fmt.Print("listen http://localhost:9947")
	err := http.ListenAndServe(":9947", nil)
	if err != nil {
		panic(err)
	}
}
