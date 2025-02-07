let socket;
let username;
let userColors = {};

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("join-btn").addEventListener("click", connect);
    document.getElementById("send-btn").addEventListener("click", sendMessage);

    document.getElementById(
        "username",
    ).addEventListener("keypress", (event) => {
        if (event.key === "Enter") connect();
    });

    document.getElementById(
        "message",
        ).addEventListener("keypress", (event) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            sendMessage();
        }
    });
});

function connect() {
    username = document.getElementById("username").value.trim();
    if (!username) {
        alert("Enter a username!");
        return;
    }

    document.getElementById("username-modal").style.display = "none";
    socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
        console.log("Connected to WebSocket server");
        socket.send(JSON.stringify({ username }));
    };

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        displayMessage(msg.username, msg.text);
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
        console.warn("WebSocket closed. Reconnecting...");
        setTimeout(connect, 3000);
    };
}

function sendMessage() {
    const messageInput = document.getElementById("message");
    const message = messageInput.value.trim();

    if (message && socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ username, text: message }));
        messageInput.value = "";
    }
}

function displayMessage(sender, text) {
    const msgContainer = document.getElementById("messages");

    if (!userColors[sender]) {
        userColors[sender] = `linear-gradient(to right, #${
            Math.random().toString(16).slice(2, 8)
        }, #${Math.random().toString(16).slice(2, 8)})`;
    }

    const msgWrapper = document.createElement("div");
    msgWrapper.classList.add("message-container", sender === username ? "self" : "other");

    const avatar = document.createElement("div");
    avatar.classList.add("avatar");
    avatar.innerText = sender[0].toUpperCase();
    avatar.style.background = userColors[sender];

    const textDiv = document.createElement("div");
    textDiv.classList.add("message", sender === username ? "self" : "other");

    const nameDiv = document.createElement("span");
    nameDiv.classList.add("username");
    nameDiv.innerText = sender;

    textDiv.appendChild(nameDiv);
    textDiv.appendChild(document.createTextNode(text));

    msgWrapper.appendChild(avatar);
    msgWrapper.appendChild(textDiv);
    msgContainer.appendChild(msgWrapper);
    msgContainer.scrollTop = msgContainer.scrollHeight;
}
