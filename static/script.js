
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

async function createChatroom(lobbyEP, conn) {

    let roomname = document.getElementById("roomname").value;
    
    roomname.trim();

    if (roomname != "") {
        
        const lobbyQuery = `http://${lobbyEP}?roomname=*&roompath=*`
        const response = await fetch(lobbyQuery);
        console.log(`${response.status} ${response.statusText}`);


        const existingrooms = await response.json();
        const existingnames = Object.keys(existingrooms);
        const existingpaths = Object.values(existingrooms);
        console.log(existingrooms)

        if (!response.ok || !existingnames.includes(roomname)) {

            console.log(`${roomname} room name available to register in lobby`);
            let roompath;
            do {
                roompath = getRandomString();
            } while (existingpaths.includes(roompath));

            const message = {data: `{"chatroom": {"name": "${roomname}", "path": "${roompath}"}}`}
            await addRoomToLobbyDB(lobbyEP, message);
            conn.send(`${message.data}`);
            roomname = "";

        } else {
            const message = `${roomname} room name already registered in lobby`
            console.log(message);
            roomname = ""
            alert(message);
        }

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

async function getLobbyChatrooms(lobbyEP) {

    const lobbyQuery = `http://${lobbyEP}?roomname=*&roompath=*`;
    const response = await fetch(lobbyQuery);
    console.log(`${response.status} ${response.statusText}`);
    const payload = response.json()

    return payload

}

async function addRoomToLobbyDB(lobbyEP, message) {
    const roominfo = JSON.parse(message.data);
    const roompath = roominfo.chatroom.path;
    const roomname = roominfo.chatroom.name;
    const lobbyQuery = `http://${lobbyEP}?roomname=${roomname}&roompath=${roompath}`;
    const response = await fetch(lobbyQuery, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
    });
    console.log("add room to lobby request - "+response.status)
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

async function changeName(conn, usersEP, messagesEP, roompath) {

    const outputName = document.getElementById('sender');
    const newName = await getNameInput(conn, usersEP, roompath, outputName.value);

    if (newName != 'undefined') {
        const path = roompath.replace("/chatroom/","");
        const userQuery = `http://${usersEP}?displayname=${outputName.value}&roompath=${path}&newname=${newName}`;
        await fetch(userQuery, {
            method: "PUT",
            headers: {
                'Content-Type': 'application/json'
            },
        });
        const msg = `${outputName.value} changed their name to ${newName}`;
        await logMessageToDB(messagesEP, roompath, msg);
        conn.send(msg);
        outputName.value = newName;
    } else {
        let nameInput = document.getElementById('nameInput');
        nameInput.value = outputName.value
    }

  }

async function logMessageToDB(messagesEP, roompath, messageString) {
    const path = roompath.replace("/chatroom/","");
    const messageQuery = `http://${messagesEP}?roompath=${path}&chatmessage=${messageString}`;
    await fetch(messageQuery, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
    });
}

async function sendMessage(messagesEP, conn, message, sender, enteredName, roompath) {

    if (sender.value != enteredName.value) {
        enteredName.value = sender.value
    }
    if (message != null) {
        const messageString = `${sender.value}: ${message.value}`;
        conn.send(messageString);
        message.value = "";
        await logMessageToDB(messagesEP, roompath, messageString)
    }

}

async function roomEntranceMessage(messagesEP, conn, displayname, roompath) {
    
    const messageString = `${displayname} has entered the chat`
    await logMessageToDB(messagesEP, roompath, messageString)
    conn.send(messageString);

}

function receiveMessage(message) {

    let chatbox = document.getElementById("chatmessages");
    console.log(message);
    const newMessage = message.data;
    chatbox.value += "\n" + newMessage;
    chatbox.scrollTop = chatbox.scrollHeight;

}

async function populateMessages(messagesEP, roompath) {

    const path = roompath.replace("/chatroom/","");
    const messageQuery = `http://${messagesEP}?roompath=${path}`;
    response = await fetch(messageQuery);
    const messages = await response.json();
    console.log(messages);
    messages.forEach(messagestr => {
        const message = {data: messagestr};
        receiveMessage(message);
    });

}

window.onload = async () => {

    if (window["WebSocket"]) {
        console.log("browser websocket support found");
        
        let pageHost = window.location.host;
        let pagePath = window.location.pathname === undefined ? "/" : window.location.pathname;
        let socketURL;

        const lobbyEP = pageHost + "/api/lobby";
        const messagesEP = pageHost + "/api/messages";
        const usersEP = pageHost + "/api/users";

        if (pagePath[pagePath.length-1] === "/") {
            socketURL = "ws://" + pageHost + "/ws_lobby";
        } else {
            socketURL = "ws://" + pageHost + "/ws_chatroom" + pagePath;
        }

        let chatmessage = document.getElementById("chatroom-message");
        let createroom = document.getElementById("chatroom-create");
        let nameInput = document.getElementById("nameInput");
        let nameInputButton = document.getElementById("nameInputButton");
        let displayname = document.getElementById("sender");
        const newmessage = document.getElementById("message");


        if (chatmessage) {
            await populateMessages(messagesEP, pagePath);
            const defaultName = await generateAnon(usersEP, pagePath);
            nameInput.value = defaultName;
            displayname.value = defaultName;
            let conn = new WebSocket(socketURL);
            conn.onopen = () => {
                const payload = defaultName;
                conn.send(payload);
            };
            roomEntranceMessage(messagesEP, conn, displayname.value, pagePath);

            window.onunload = async () => {
                await removeDisplayNameFromRoom(pagePath, usersEP, displayname.value);
            };
            
            nameInputButton.onclick = async () => {
                await changeName(conn, usersEP, messagesEP, pagePath);
            };
            chatmessage.onsubmit = (event) => {  
                event.preventDefault();
                sendMessage(messagesEP, conn, newmessage, displayname, nameInput, pagePath);
            };
            conn.onmessage = (message) => {
                if (message.data != "client disconnect") {
                    receiveMessage(message);
                }
            };
        }

        if (createroom) {
            let conn = new WebSocket(socketURL);
            conn.onopen = () => {
                const payload = 'ClientNameCookiePlaceholder';
                conn.send(payload);
            };
            console.log(conn);
            const rooms = await getLobbyChatrooms(lobbyEP);
            console.log(rooms);
            for (const k in rooms) {
                let roompath = rooms[k].replace('"', '').replace('"', '');
                const room = {data: `{"chatroom": {"name": "${k}", "path": "${roompath}"}}`};
                updateRoomList(room);
            }
            
            createroom.onsubmit = (event) => {
                console.log('create room on submit')
                event.preventDefault();
                createChatroom(lobbyEP, conn);
            };
            conn.onmessage = (message) => {
                console.log(message)
                if (message.data != "client disconnect") {
                    updateRoomList(message);
                }
            };
        }


    } else {
        alert("Websockets not supported by browser!");
    }
};

