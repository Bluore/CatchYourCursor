const Host = window.location.host

var WSClient = new WebSocket(`ws://${Host}/api/cursor`)

WSClient.onmessage = function (event) {
    console.log("recieve: ",event.data)
    const data = JSON.parse(event.data)
    switch (data.status){
        case 1:
            MoveTheCursor(data.data);
            break;
    }
}

function MoveTheCursor(data) {
    console.log(data)
    console.log(`move the cursor to (${data.x},${data.y})`)
    const cursor = document.getElementById("test")
    cursor.style.left = `${data.x}px`
    cursor.style.top = `${data.y}px`
}