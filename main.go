package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
  <head>
    <title>websocket</title>
  </head>

  <body>
    <samp id="lines"></sampo>
    <script type="text/javascript">
(function(){
    function connect(uri) {
        var socket = new WebSocket(uri);
        var elem = document.getElementById("lines");
        socket.addEventListener("open", function (e) {
            console.log("open websocket");
            elem.innerHTML = "";
        });

        socket.addEventListener("close", function (e) {
            console.log("close websocket");
            setTimeout(reconnect, 1000);
        });

        socket.addEventListener("error", function (e) {
            console.log("error!!", e);
        });

        socket.addEventListener("message", function (e) {
            var data = JSON.parse(e.data);
            var line = document.createElement("div");
            line.innerText = data.text;
            elem.appendChild(line);
        });
    }
    var uri = "ws://" + location.host + "/follow";
    connect(uri);
    function reconnect() {
        connect(uri);
    }
})();
    </script>
  </body>
</html>
`)
}

func main() {
	t, _ := NewTailReader(os.Stdin)

	http.Handle("/follow", websocket.Handler(t.FollowHandler))
	http.HandleFunc("/", IndexHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
