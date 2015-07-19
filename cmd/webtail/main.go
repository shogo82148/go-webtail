package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/shogo82148/go-webtail"
	"golang.org/x/net/websocket"
)

func IndexHandler(prefix string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, indexHTML, prefix)
	}
}

func main() {
	var port int
	var prefix string
	var bufferLines int
	var playBackLines int

	flag.IntVar(&port, "port", 8080, "listen port(default: 8080)")
	flag.IntVar(&port, "p", 8080, "listen port")
	flag.StringVar(&prefix, "prefix", "", "prefix of url")
	flag.IntVar(&bufferLines, "buffer", 10240, "buffering lines(default: 10240)")
	flag.IntVar(&playBackLines, "playback", 10, "auto play back lines(default: 10)")
	flag.Parse()

	if len(prefix) > 1 && prefix[0] != '/' {
		prefix = "/" + prefix
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}

	for _, file := range args {
		if file == "-" {
			// tail stdin
			t, _ := webtail.NewTailReader(os.Stdin)
			t.BufferLines = bufferLines
			t.PlayBackLines = playBackLines
			http.Handle(prefix+"/tail", http.StripPrefix(prefix, http.HandlerFunc(t.TailHandler)))
			http.Handle(prefix+"/follow", http.StripPrefix(prefix, websocket.Handler(t.FollowHandler)))
			http.HandleFunc(prefix+"/", IndexHandler(prefix))
		} else {
			// tail file
			t, _ := webtail.NewTailFile(file)
			t.BufferLines = bufferLines
			t.PlayBackLines = playBackLines
			basename := path.Base(file)
			http.Handle(prefix+"/tail", http.StripPrefix(prefix, http.HandlerFunc(t.TailHandler)))
			http.Handle(prefix+"/"+basename+"/follow", http.StripPrefix(prefix+"/"+basename, websocket.Handler(t.FollowHandler)))
			http.HandleFunc(prefix+"/"+basename, IndexHandler(prefix+"/"+basename))
		}
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
