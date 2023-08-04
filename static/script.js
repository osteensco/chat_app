
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

    let roomname = document.getElementById("roomname").value;
    
    roomname.trim();

    if (roomname != "") {
        let roompath = getRandomString();
        conn.send(`{"chatroom": {"name": "${roomname}", "path": "${roompath}"}}`)
        roomname.value = ""
    }
}

function updateRoomList(message) {

    let roomlist = document.getElementById("room-list");
    let newroom = document.createElement("option");
    
    console.log(message);
    if (message.data === "client disconnect") {
        return
    }

    let roominfo = JSON.parse(message.data);
    newroom.value = roominfo.chatroom.path;
    newroom.text = roominfo.chatroom.name;
    roomlist.appendChild(newroom);
}

function navToChatroom() {
    const room = document.getElementById("room-list");
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
    const newmessage = document.getElementById("message");
    const sender = document.getElementById("sender");
    if (newmessage != null) {
        conn.send(`${sender.value}: ${newmessage.value}`);
        newmessage.value = "";
    }
}

function receiveMessage(message) {
    let chatbox = document.getElementById("chatmessages");
    console.log(message);
    const newMessage = message.data;
    chatbox.value += "\n" + newMessage;
    chatbox.scrollTop = chatbox.scrollHeight;

}

window.onload = function () {

    if (window["WebSocket"]) {
        console.log("browser websocket support found");
        
        let pageHost = window.location.host;
        let pagePath = window.location.pathname === undefined ? "/" : window.location.pathname;
        let socketURL;


        if (pagePath[pagePath.length-1] === "/") {
            socketURL = "ws://" + pageHost + "/ws";
        } else {
            socketURL = "ws://" + pageHost + "/ws" + pagePath;
        }
        
        let conn = new WebSocket(socketURL);

        let chatmessage = document.getElementById("chatroom-message");
        let createroom = document.getElementById("chatroom-create");

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