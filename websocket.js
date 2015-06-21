var uri = "ws://" + location.host + "/follow";
var socket = new WebSocket(uri);

socket.addEventListener("open", function (e) {
    document.getElementById("serverStatus").innerHTML =
        'WebSocket Status:: Socket Open';
});

socket.addEventListener("error", function (e) {
    document.getElementById("serverStatus").innerHTML =
        'WebSocket Status:: Socket Error';
});

socket.addEventListener("close", function (e) {
    document.getElementById("serverStatus").innerHTML =
        'WebSocket Status:: Socket Close';
});

socket.addEventListener("message", function (e) {
    var data = JSON.parse(e.data);
    document.getElementById("messages").innerHTML += "<div>" + data.text + "</div>";
});
