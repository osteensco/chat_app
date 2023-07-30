
function createChatroom() {

    var roomlist = document.getElementById("room-list");
    var roomname = document.getElementById("roomname").value;
    
    roomname.trim();

    if (roomname != "") {
        newroom = document.createElement("option");
        newroom.value = "chatroom.html";//change to a random value
        //send room name to server
        //wait for confirmation room has been created on the server
        newroom.text = roomname;
        
        roomlist.appendChild(newroom);

        roomname.value = ""
    }
}

function navToChatroom() {
    var room = document.getElementById("room-list").value;
    if (room) {
        // send request to server with room name path
        // 
    }

}

function changeName() {
    const nameInput = document.getElementById('nameInput').value;
    const outputName = document.getElementById('sender');
    outputName.value = nameInput || 'Anonymous';
  }

function sendMessage(conn) {
    var newmessage = document.getElementById("message");
    var sender = document.getElementById("sender");
    if (newmessage != null) {
        conn.send(`${sender.value}: ${newmessage.value}`);
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

        document.getElementById("chatroom-create").onsubmit = (event) => {
            event.preventDefault();
            createChatroom();
        };
        
        conn.onmessage = (message) => {
            receiveMessage(message)
        };

    } else {
        alert("Websockets not supported by browser!");
    }
};