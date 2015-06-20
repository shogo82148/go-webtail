package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"github.com/mattn/go-pubsub"
	"golang.org/x/net/websocket"
)

type Line struct {
	Num int
	Str string
}

var ps *pubsub.PubSub

func echoHandler(ws *websocket.Conn) {
	ch := make(chan Line)
	sub := func(l Line) { ch <- l }
	ps.Sub(sub)
	defer ps.Leave(sub)

	for {
		line := <-ch
		if err := websocket.JSON.Send(ws, line); err != nil {
			log.Print(err)
			break
		}
	}
}

func tail() {
	buf := bufio.NewReader(os.Stdin)
	num := 1
	for {
		b, _, _ := buf.ReadLine()
		if len(b) > 0 {
			ps.Pub(Line{Num: num, Str: string(b)})
			num++
		}
	}
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	ps = pubsub.New()
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	tail()
}
