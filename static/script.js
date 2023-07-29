


function sendMessage() {
    var newmessage = document.getElementById("message");
    if (newmessage != null) {
        console.log(newmessage);
        conn = new WebSocket("ws://" + document.location.host + "/");
    }
    return false;
}



window.onload = function () {

    document.getElementById("chatroom-message").onsubmit = sendMessage;

    if (window["WebSocket"]) {
        console.log("browser websocket support found");
    } else {
        alert("Websockets not supported by browser!");
    }
};