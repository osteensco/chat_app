


function sendMessage(conn) {
    var newmessage = document.getElementById("message");
    if (newmessage != null) {
        conn.send(newmessage.value);
        newmessage.value = "";
    }
    return false;
}

function receiveMessage(message) {
    var chatbox = document.getElementById("chatmessages");
    console.log(message);
    var newMessage = message.data;
    chatbox.value += "\n" + newMessage;
    chatbox.scrollTop = chatbox.scrollHeight;

}

window.onload = function () {

    if (window["WebSocket"]) {
        console.log("browser websocket support found");
        conn = new WebSocket("ws://" + document.location.host + "/ws");

        document.getElementById("chatroom-message").onsubmit = (event) => {
            event.preventDefault();
            sendMessage(conn);
        };
        
        conn.onmessage = (message) => {
            receiveMessage(message)
        };

    } else {
        alert("Websockets not supported by browser!");
    }
};