
function getRandomString() {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@$';
    const minStringLength = 16;
    const maxStringLength = 24;
    const randomStringLength = Math.floor(Math.random() * (maxStringLength - minStringLength + 1)) + minStringLength;
  
    let randomString = '';
    for (let i = 0; i < randomStringLength; i++) {
      const randomIndex = Math.floor(Math.random() * characters.length);
      randomString += characters[randomIndex];
    }
  
    return randomString;
  }

function createChatroom(conn) {

    var roomname = document.getElementById("roomname").value;
    
    roomname.trim();

    if (roomname != "") {
        var roompath = getRandomString();
        conn.send(`{"chatroom": {"name": "${roomname}", "path": "${roompath}"}}`)
        roomname.value = ""
    }
}

function updateRoomList(message) {

    var roomlist = document.getElementById("room-list");

    var newroom = document.createElement("option");
    console.log(message);
    var roominfo = JSON.parse(message.data);
    newroom.value = roominfo.chatroom.path;
    newroom.text = roominfo.chatroom.name;
    roomlist.appendChild(newroom);
}

function navToChatroom() {
    var room = document.getElementById("room-list");
    if (room) {
        window.location.href = `/chatroom/${room.value}`
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
        
        var pageHost = window.location.host;
        var pagePath = window.location.pathname === undefined ? "/" : window.location.pathname;
        var socketPath;
        console.log(pagePath)
        console.log(typeof(pagePath))
        if (pagePath[pagePath.length-1] === "/") {
            socketPath = "ws"

        } else {
            socketPath = "/ws"
        }
        var socketURL = "ws://" + pageHost + pagePath + socketPath;
        console.log(socketURL)
        var conn = new WebSocket(socketURL);

        var chatmessage = document.getElementById("chatroom-message");
        var createroom = document.getElementById("chatroom-create");

        if (chatmessage) {
            chatmessage.onsubmit = (event) => {  
                event.preventDefault();
                sendMessage(conn);
            };
            conn.onmessage = (message) => {
                receiveMessage(message);
            };
        }

        if (createroom) {
            createroom.onsubmit = (event) => {
                event.preventDefault();
                createChatroom(conn);
            };
            conn.onmessage = (message) => {
                updateRoomList(message);
            }
        }
        


    } else {
        alert("Websockets not supported by browser!");
    }
};