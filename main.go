package main

import (
	"fmt"
	"net/http"
)

const port = 9947

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/api/cursor", func(w http.ResponseWriter, r *http.Request) { CursorWebsocketHandler(w, r, hub) })

	http.Handle("/", http.FileServer(http.Dir("./public")))

	fmt.Printf("listen http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
