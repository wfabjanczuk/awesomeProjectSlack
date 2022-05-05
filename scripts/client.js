const socket = new WebSocket("ws://localhost:8080/ws");

socket.send(JSON.stringify({
    action: "broadcast",
    message: "Hello world"
}));
