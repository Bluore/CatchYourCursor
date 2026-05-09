const Host = window.location.host

var WSClient = new WebSocket(`ws://${Host}/api/cursor`)
clientID = ""

WSClient.onmessage = function (event) {
    console.log("recieve: ",event.data)
    const data = JSON.parse(event.data)
    switch (data.status){
        case 1:
            MoveTheCursor(data.data);
            break;
        case 6:
            UpdateClientID(data.data);
            break;
    }
}

function UpdateClientID(data) {
    clientID = data.id
}

function MoveTheCursor(data) {
    console.log(data)
    console.log(`move the cursor to (${data.x},${data.y})`)
    const cursor = document.getElementById("test")
    cursor.style.left = `${data.x}px`
    cursor.style.top = `${data.y}px`
}

mouseMoveLastTime = 0
document.addEventListener("mousemove",(e)=>{
    const mouseMoveNow = Date.now();
    if (mouseMoveNow - mouseMoveLastTime < 100){
        return;
    }
    mouseMoveLastTime = mouseMoveNow;
    console.log(`mouse is in (${e.x},${e.y})`);
    
    const sendData = {
        status: 1,
        id: clientID,
        data: {
            id: clientID,
            x: e.x,
            y: e.y,
        }
    };
    sendMsgToWSClient(sendData);
})

async function sendMsgToWSClient(data){
    const sendMsg = JSON.stringify(data)
    WSClient.send(sendMsg)
    console.log(`already send msg: ${sendMsg}`)
}