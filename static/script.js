
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

async function generateAnon(usersEP, pagePath) {

    let anon = "Anonymous"
    const min = 0;
    const max = 999999;
    let randomNumber = Math.floor(Math.random() * (max - min + 1)) + min;
    randomNumber = randomNumber.toString().padStart(6, '0');
    anon = anon.concat(randomNumber);
    
    anon = await checkDisplayNameAvailability(() => {
        generateAnon(usersEP, pagePath);
    }, pagePath, usersEP, anon);

    addDisplayNameToRoom(pagePath, usersEP, anon);

    return anon

    }

async function checkDisplayNameAvailability(callback, pagePath, usersEP, displayname) {

    const path = pagePath.replace("/chatroom/","");
    const userQuery = `http://${usersEP}?displayname=${displayname}&roompath=${path}`;

    const response = await fetch(userQuery);
    console.log(`${response.status} ${response.statusText}`);

    if (!response.ok) {
        console.log(`${displayname} display name available to register in ${path}`);
        return displayname;
    } else {
        console.log(`${displayname} display name already registered in ${path}`);
        callback();
    }

}

async function addDisplayNameToRoom(pagePath, usersEP, displayname) {

    const path = pagePath.replace("/chatroom/","");
    const userQuery = `http://${usersEP}?displayname=${displayname}&roompath=${path}`;

    const response = await fetch(userQuery, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
    });

    console.log(`${response.status} ${response.statusText}`);

}

async function removeDisplayNameFromRoom(pagePath, usersEP, displayname) {

    const path = pagePath.replace("/chatroom/","");
    const userQuery = `http://${usersEP}?displayname=${displayname}&roompath=${path}`;

    const response = await fetch(userQuery, {
        method: "DELETE",
        headers: {
            'Content-Type': 'application/json'
        },
    });

    console.log(`${response.status} ${response.statusText}`);

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



async function getNameInput(conn, usersEP, pagePath) {
    let nameInput = document.getElementById('nameInput');
    if (nameInput.value != '') {    
        nameInput.value = await checkDisplayNameAvailability(() => {
            alert(`name ${nameInput.value} already taken, please choose another`);
        }, pagePath, usersEP, nameInput.value);
    }

    return nameInput.value
}

async function changeName(conn, usersEP, pagePath) {

    const outputName = document.getElementById('sender');
    const newName = await getNameInput(conn, usersEP, pagePath, outputName.value);

    if (newName != 'undefined') {
        const path = pagePath.replace("/chatroom/","");
        const userQuery = `http://${usersEP}?displayname=${outputName.value}&roompath=${path}&newname=${newName}`;
        await fetch(userQuery, {
            method: "PUT",
            headers: {
                'Content-Type': 'application/json'
            },
        });
        conn.send(`${outputName.value} changed their name to ${newName}`);
        outputName.value = newName;
    } else {
        let nameInput = document.getElementById('nameInput');
        nameInput.value = outputName.value
    }


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
    // /api/messages
}

function receiveMessage(message) {
    let chatbox = document.getElementById("chatmessages");
    console.log(message);
    const newMessage = message.data;
    chatbox.value += "\n" + newMessage;
    chatbox.scrollTop = chatbox.scrollHeight;

}

window.onload = async function () {

    if (window["WebSocket"]) {
        console.log("browser websocket support found");
        
        let pageHost = window.location.host;
        let pagePath = window.location.pathname === undefined ? "/" : window.location.pathname;
        let socketURL;

        // TODO
        // These need to be passed in to a function at somepoint likely, need to finalized API call flow from frontend
        const lobbyEP = pageHost + "/api/lobby";
        const messagesEP = pageHost + "/api/messages";
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
            // /api/messages
            const defaultName = await generateAnon(usersEP, pagePath);
            nameInput.value = defaultName;
            displayname.value = defaultName;

            window.onunload = async () => {
                await removeDisplayNameFromRoom(pagePath, usersEP, displayname.value);
            };
            
            conn.onopen = () => {
                conn.send(`${displayname.value} has entered the chat`);
            };
            nameInputButton.onclick = async () => {
                await changeName(conn, usersEP, pagePath);
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

