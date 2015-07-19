package main

const indexHTML = `<!DOCTYPE html>
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
    var socket;
    function connect(uri) {
        socket = new WebSocket(uri);
        var elem = document.getElementById("lines");

       function addNewLine(newLine) {
            var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
            var windowHeight = document.documentElement.clientHeight;
            var needScroll = scrollTop + windowHeight >= elem.scrollTop + elem.scrollHeight;
            elem.appendChild(newLine);
            if (needScroll) window.scroll(0, elem.scrollTop + elem.scrollHeight);
       }

        socket.addEventListener("open", function (e) {
            var line = document.createElement("div");
            line.className = "system-message";
            line.innerText = "connection connected";
            addNewLine(line);
            reconnecting = false;
        });

        socket.addEventListener("close", function (e) {
            socket = undefined;
            if (!reconnecting) {
                var line = document.createElement("div");
                line.className = "system-message";
                line.innerText = "connection closed";
                addNewLine(line);
            }
            reconnecting = true;
            setTimeout(reconnect, 10000);
        });

        socket.addEventListener("error", function (e) {
            if (!reconnecting) {
                var line = document.createElement("div");
                line.className = "system-message";
                line.innerText = "connection error";
                addNewLine(line);
            }
        });

        socket.addEventListener("message", function (e) {
            var data = JSON.parse(e.data);
            var line = document.createElement("div");
            line.innerText = data.text;
            addNewLine(line);
        });
    }
    var uri = "ws://" + location.host + "%s/follow";
    connect(uri);
    function reconnect() {
        connect(uri);
    }

    setInterval(function() {
        if (socket) socket.send('ping');
    }, 20*1000);

})();
    </script>
  </body>
</html>
`
