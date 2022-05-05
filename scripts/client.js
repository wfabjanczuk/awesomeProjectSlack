const socket = new WebSocket("ws://localhost:8080/ws");
// Server: Client connected to the endpoint

socket.send(JSON.stringify({
    action: "broadcast",
    message: "Hello world"
}));
// Server: New message from client, action: "broadcast", message: "Hello world"
// Client: {"message":"Hello world","status":"OK"}

socket.send(JSON.stringify({
    action: "broadcast",
    message: "Second message"
}));
// Server: New message from client, action: "broadcast", message: "Second message"
// Client: {"message":"Second message","status":"OK"}

socket.send(JSON.stringify({
    action: "create",
    message: "public"
}));
// Server: New message from client, action: "create", message: "public"
// Client: {"message":"Channel with name: \"public\" already exists!","status":"Error"}

socket.send(JSON.stringify({
    action: "create",
    message: "private"
}));
// Server: New message from client, action: "create", message: "private"
// Client: {"message":"Channel with name: \"private\" successfully created.","status":"OK"}