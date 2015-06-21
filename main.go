package main

import (
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	t, _ := NewTailReader(os.Stdin)

	http.Handle("/follow", websocket.Handler(t.FollowHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
