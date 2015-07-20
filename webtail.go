package webtail

import (
	"container/list"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mattn/go-pubsub"
	"github.com/shogo82148/go-tail"
	"golang.org/x/net/websocket"
)

type Line struct {
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
	Number int64     `json:"number"`
}

type Tail struct {
	// BufferLines is buffer size for play back
	BufferLines int

	// PlayBackLines is the number of lines for auto play back
	PlayBackLines int

	ps     *pubsub.PubSub
	t      *tail.Tail
	mu     sync.RWMutex
	lines  *list.List
	number int64
}

func NewTail(t *tail.Tail) (*Tail, error) {
	myt := &Tail{
		ps:    pubsub.New(),
		t:     t,
		lines: list.New(),
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
			t.addNewLine(line)
		case err := <-t.t.Errors:
			log.Print("Error: ", err)
		}
	}
}

func (t *Tail) addNewLine(newline *tail.Line) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.number++
	line := &Line{
		Text:   newline.Text,
		Time:   newline.Time,
		Number: t.number,
	}
	t.ps.Pub(line)

	if t.BufferLines > 0 {
		t.lines.PushBack(line)
		for t.lines.Len() > t.BufferLines {
			t.lines.Remove(t.lines.Front())
		}
	}
}

func (t *Tail) TailHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	strLines := values.Get("lines")
	numLines := t.PlayBackLines
	if strLines != "" {
		numLines, _ = strconv.Atoi(strLines)
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	e := t.lines.Back()
	for i := 1; e != nil && i < numLines; e = e.Prev() {
		i++
	}
	if numLines <= 0 || e == nil {
		e = t.lines.Front()
	}

	lines := []*Line{}
	for ; e != nil; e = e.Next() {
		line, ok := e.Value.(*Line)
		if !ok {
			continue
		}
		lines = append(lines, line)
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(lines)
}

func (t *Tail) FollowHandler(ws *websocket.Conn) {
	// start subscribe
	ch := make(chan *Line)
	sub := func(l *Line) { ch <- l }
	t.ps.Sub(sub)
	defer t.ps.Leave(sub)

	go func() {
		for {
			var message string
			err := websocket.Message.Receive(ws, &message) // ignore message
			if err != nil {
				break
			}
		}
	}()

	// wait new lines
	for {
		line := <-ch
		if err := websocket.JSON.Send(ws, line); err != nil {
			log.Print(err)
			break
		}
	}
}
