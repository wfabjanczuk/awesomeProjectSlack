const socket = new WebSocket("ws://localhost:8080/ws");

socket.send(JSON.stringify({
    action: "broadcast",
    message: "Hello world"
}));

socket.send(JSON.stringify({
    action: "broadcast",
    message: "Second message"
}));