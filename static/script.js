
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

function generateAnon(usersEP, pagePath) {

    let anon = "Anonymous"
    const min = 0;
    const max = 999999;
    let randomNumber = Math.floor(Math.random() * (max - min + 1)) + min;
    randomNumber = randomNumber.toString().padStart(6, '0');
    anon = anon.concat(randomNumber);
    
    anon = checkDisplayNameAvailability(() => {
        generateAnon(usersEP, pagePath);
    }, pagePath, usersEP, anon);

    return anon

    }

function checkDisplayNameAvailability(callback, pagePath, usersEP, displayname) {

    const path = pagePath.replace("/chatroom/","");
    const userQuery = `http://${usersEP}?displayname=${displayname}&roompath=${path}`;
   
    fetch(userQuery)
        .then(response => {
            if (!response.ok) {
                console.log(`${response.status} ${response.statusText}`)
                console.log(`${displayname} display name available to register in ${path}`)
                // TODO 
                // add aditional logic to do PUT request for name change events
                // CREATE record in cockroachDB
                fetch(userQuery, {method: "POST", 
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ 
                        chatroom_path: path,  
                        display_name: displayname
                    }) 
                })
                    .then(response => {
                        console.log(`added displayname ${displayname} to set in redis`)
                        console.log(response)
                    })    
            
            } else {
                callback();
                return response.json();
            }
        })
        .then(data => {
            console.log(`data: ${data}`)
        })


    return displayname
}

function createChatroom(conn) {

    let roomname = document.getElementById("roomname").value;
    
    roomname.trim();

    if (roomname != "") {
        let roompath = getRandomString();
        //TODO
        // READ record in redis to make sure roomname and/or roompath doesn't already exist
        // /api/lobby
        conn.send(`{"chatroom": {"name": "${roomname}", "path": "${roompath}"}}`)
        // CREATE record in redis and cockroachDB
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

function displayInputMessage(nameInput) {
    alert(`name ${nameInput} already taken, please choose another`);
}

function getNameInput(conn, usersEP, pagePath) {
    let nameInput = document.getElementById('nameInput');
    if (nameInput.value != '') {    
        nameInput.value = checkDisplayNameAvailability(() => {
            displayInputMessage(nameInput.value);
        }, pagePath, usersEP, nameInput.value);
    }

    return nameInput.value
}

function changeName(conn, usersEP, pagePath) {
    const outputName = document.getElementById('sender');
    const newName = getNameInput(conn, usersEP, pagePath);

    if (outputName.value != newName) {
        conn.send(`${outputName.value} changed their name to ${newName}`)
    }

    outputName.value = newName;
  }

function sendMessage(conn, message, sender, enteredName) {
    if (sender.value != enteredName.value) {
        enteredName.value = sender.value
    }
    if (message != null) {
        conn.send(`${sender.value}: ${message.value}`);
        message.value = "";
    }
    // TODO
    // CREATE or UPDATE record in redis and cockroachDB
    // /api/chatrooms
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

        // TODO
        // These need to be passed in to a function at somepoint likely, need to finalized API call flow from frontend
        const lobbyEP = pageHost + "/api/lobby";
        const chatroomsEP = pageHost + "/api/chatrooms";
        const usersEP = pageHost + "/api/users";

        if (pagePath[pagePath.length-1] === "/") {
            socketURL = "ws://" + pageHost + "/ws_lobby";
        } else {
            socketURL = "ws://" + pageHost + "/ws_chatroom" + pagePath;
        }
        
        let conn = new WebSocket(socketURL);

        let chatmessage = document.getElementById("chatroom-message");
        let createroom = document.getElementById("chatroom-create");
        let nameInput = document.getElementById("nameInput");
        let nameInputButton = document.getElementById("nameInputButton");
        let displayname = document.getElementById("sender");
        const newmessage = document.getElementById("message");
        

        if (chatmessage) {
            // TODO
            // READ redis and display recent chat messages (last 10? 20?)
            // /api/chatrooms
            const defaultName = generateAnon(usersEP, pagePath)
            nameInput.value = defaultName
            displayname.value = defaultName
            
            conn.onopen = () => {
                conn.send(`${displayname.value} has entered the chat`)
            };
            nameInputButton.onclick = () => {
                changeName(conn, usersEP, pagePath);
            };
            chatmessage.onsubmit = (event) => {  
                event.preventDefault();
                sendMessage(conn, newmessage, displayname, nameInput);
            };
            conn.onmessage = (message) => {
                if (message.data != "client disconnect") {
                    receiveMessage(message);
                }
            };
        }

        if (createroom) {
            // TODO
            // READ redis and display current available chatrooms
            // /api/lobby
            createroom.onsubmit = (event) => {
                event.preventDefault();
                createChatroom(conn);
            };
            conn.onmessage = (message) => {
                if (message.data != "client disconnect") {
                    updateRoomList(message);
                }
            };
        }


    } else {
        alert("Websockets not supported by browser!");
    }
};