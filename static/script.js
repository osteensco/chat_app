


function sendMessage(conn) {
    var newmessage = document.getElementById("message");
    if (newmessage != null) {
        console.log(newmessage);

        conn.send(newmessage.value);
    }
    console.log("error sending message")
    return false;
}



window.onload = function () {

    if (window["WebSocket"]) {
        console.log("browser websocket support found");
        conn = new WebSocket("ws://" + document.location.host + "/");

        document.getElementById("chatroom-message").onsubmit = () => {console.log("error"); sendMessage(conn);};
        
        conn.onmessage = (message) => {
            console.log(message);
        }

    } else {
        alert("Websockets not supported by browser!");
    }
};