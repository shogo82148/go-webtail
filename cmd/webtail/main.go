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
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <title>go-webtail</title>
    <style type="text/css">
    .system-message { color: lightsteelblue; }
    </style>
  </head>

  <body>
    <pre id="lines"></pre>
    <script type="text/javascript">
(function(){
    var reconnecting = true;
    function connect(uri) {
        var socket = new WebSocket(uri);
        var elem = document.getElementById("lines");
        socket.addEventListener("open", function (e) {
            var line = document.createElement("div");
            line.className = "system-message";
            line.innerText = "connection connected";
            elem.appendChild(line);
            reconnecting = false;
        });

        socket.addEventListener("close", function (e) {
            if (!reconnecting) {
                var line = document.createElement("div");
                line.className = "system-message";
                line.innerText = "connection closed";
                elem.appendChild(line);
            }
            reconnecting = true;
            setTimeout(reconnect, 1000);
        });

        socket.addEventListener("error", function (e) {
            if (!reconnecting) {
                var line = document.createElement("div");
                line.className = "system-message";
                line.innerText = "connection error";
                elem.appendChild(line);
            }
        });

        socket.addEventListener("message", function (e) {
            var data = JSON.parse(e.data);
            var line = document.createElement("div");
            line.innerText = data.text;
            elem.appendChild(line);
        });
    }
    var uri = "ws://" + location.host + "%s/follow";
    connect(uri);
    function reconnect() {
        connect(uri);
    }
})();
    </script>
  </body>
</html>
`, prefix)
	}
}

func main() {
	var port int
	var prefix string

	flag.IntVar(&port, "port", 8080, "listen port(default: 8080)")
	flag.IntVar(&port, "p", 8080, "listen port")
	flag.StringVar(&prefix, "prefix", "", "prefix of url")
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
			http.Handle(prefix+"/follow", http.StripPrefix(prefix, websocket.Handler(t.FollowHandler)))
			http.HandleFunc(prefix+"/", IndexHandler(prefix))
		} else {
			// tail file
			t, _ := webtail.NewTailFile(file)
			basename := path.Base(file)
			http.Handle(prefix+"/"+basename+"/follow", http.StripPrefix(prefix+"/"+basename, websocket.Handler(t.FollowHandler)))
			http.HandleFunc(prefix+"/"+basename, IndexHandler(prefix+"/"+basename))
		}
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
