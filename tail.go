package main

import (
	"io"
	"log"
	"time"

	"github.com/mattn/go-pubsub"
	"github.com/shogo82148/go-tail"
	"golang.org/x/net/websocket"
)

type Line struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

type Tail struct {
	ps *pubsub.PubSub
	t  *tail.Tail
}

func NewTail(t *tail.Tail) (*Tail, error) {
	myt := &Tail{
		ps: pubsub.New(),
		t:  t,
	}
	go myt.run()
	return myt, nil
}

func NewTailFile(filename string) (*Tail, error) {
	t, err := tail.NewTailFile(filename)
	if err != nil {
		return nil, err
	}
	return NewTail(t)
}

func NewTailReader(reader io.Reader) (*Tail, error) {
	t, err := tail.NewTailReader(reader)
	if err != nil {
		return nil, err
	}
	return NewTail(t)
}

func (t *Tail) run() {
	for {
		select {
		case line := <-t.t.Lines:
			t.ps.Pub(Line{
				Text: line.Text,
				Time: line.Time,
			})
		case err := <-t.t.Errors:
			log.Print("Error: ", err)
		}
	}
}

func (t *Tail) FollowHandler(ws *websocket.Conn) {
	ch := make(chan Line)
	sub := func(l Line) { ch <- l }
	t.ps.Sub(sub)
	defer t.ps.Leave(sub)

	for {
		line := <-ch
		if err := websocket.JSON.Send(ws, line); err != nil {
			log.Print(err)
			break
		}
	}
}
