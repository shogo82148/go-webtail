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
    <div><a href="?full">show all lines in buffer</a></div>
    <pre id="lines"></pre>
    <script type="text/javascript">
(function(){
    var prefix = '%s';
    var reconnecting = true;
    var playing_back = false;
    var buf = [];
    var socket;
    var xhr;
    var needScroll;
    var showAll = window.location.search === "?full";

    var elem = document.getElementById("lines");
    function addNewLines(newLines) {
        var i;
        for (i = 0; i < newLines.length; i++) {
            elem.appendChild(newLines[i]);
        }
        if (needScroll) window.scroll(0, elem.scrollTop + elem.scrollHeight);
    }

    window.addEventListener('scroll', function (e) {
        var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
        var windowHeight = document.documentElement.clientHeight;
        needScroll = scrollTop + windowHeight >= elem.scrollTop + elem.scrollHeight;
    });

    function playBack() {
        xhr = new XMLHttpRequest();
        xhr.open("GET", prefix + "/tail" + (showAll ? "?lines=0" : ""));
        showAll = false;
        xhr.send();
        playing_back = true;
        xhr.addEventListener("load", function(ev) {
            var response = JSON.parse(xhr.responseText);
            var i;
            var lines = [];
            for (i = 0; i < response.length; i++) {
                var line = document.createElement("div");
                line.innerText = response[i].text;
                lines.push(line);
            }
            var max_number = response[response.length-1].number;

            for (i = 0; i < buf.length; i++) {
                if (buf[i].number <= max_number) continue;
                var line = document.createElement("div");
                line.innerText = buf[i].text;
                lines.push(line);
            }
            addNewLines(lines);
            playing_back = false;
        });
    }

    function connect(uri) {
        socket = new WebSocket(uri);
        buf = [];

        socket.addEventListener("open", function (e) {
            var line = document.createElement("div");
            line.className = "system-message";
            line.innerText = "connection connected";
            addNewLines([line]);
            reconnecting = false;
            playBack();
        });

        socket.addEventListener("close", function (e) {
            socket = undefined;
            if (!reconnecting) {
                var line = document.createElement("div");
                line.className = "system-message";
                line.innerText = "connection closed";
                addNewLines([line]);
            }
            reconnecting = true;
            setTimeout(reconnect, 10000);
        });

        socket.addEventListener("error", function (e) {
            if (!reconnecting) {
                var line = document.createElement("div");
                line.className = "system-message";
                line.innerText = "connection error";
                addNewLines(line);
            }
        });

        socket.addEventListener("message", function (e) {
            var data = JSON.parse(e.data);
            if (playing_back) {
                buf.push(data);
            } else {
                var line = document.createElement("div");
                line.innerText = data.text;
                addNewLines([line]);
            }
        });
    }
    var uri = "ws://" + location.host + prefix + "/follow";
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
