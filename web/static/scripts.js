let socket;
let username;

function connect() {
    username = document.getElementById("username").value;
    if (!username) {
        alert("Enter a username!");
        return;
    }

    document.getElementById("username-modal").style.display = "none";

    socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
        socket.send(JSON.stringify({ username }));
    };


    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        const msgContainer = document.getElementById("messages");

        const msgDiv = document.createElement("div");
        msgDiv.classList.add("message");
        msgDiv.classList.add(msg.username === username ? "self" : "other");

        const avatar = document.createElement("div");
        avatar.classList.add("avatar");
        msgDiv.appendChild(avatar);

        const content = document.createElement("div");

        const nameDiv = document.createElement("div");
        nameDiv.classList.add("username");
        nameDiv.innerText = msg.username;
        content.appendChild(nameDiv);

        const textDiv = document.createElement("div");
        textDiv.innerText = msg.text;
        content.appendChild(textDiv);

        msgDiv.appendChild(content);
        msgContainer.appendChild(msgDiv);
        msgContainer.scrollTop = msgContainer.scrollHeight;
    };
}

function sendMessage() {
    const messageInput = document.getElementById("message");
    const message = messageInput.value;
    if (message) {
        socket.send(JSON.stringify({ username: username, text: message }));
        messageInput.value = "";
    }
}
