const Host = window.location.host
const CursorList = []

var WSClient = new WebSocket(`ws://${Host}/api/cursor`)
clientID = ""

WSClient.onmessage = function (event) {
    console.log("recieve: ",event.data)
    const data = JSON.parse(event.data)
    switch (data.status){
        case 1:
            updateCursorList(data.data);
            break;
        case 6:
            UpdateClientID(data.data);
            break;
    }
}

WSClient.onopen = function (event) {
}

function UpdateClientID(data) {
    clientID = data.id

    // 获取历史指针
    const sendData = {
        status: 5,
        id: clientID,
        data: {
            id: clientID,
            x: e.x,
            y: e.y,
        }
    };
    sendMsgToWSClient(sendData);
}

function CreateDOMCursor(data){
    const parent = document.getElementById("cursor-list");
    const child = document.createElement("div");

    child.className = "cursor-item";
    child.id = `cursor${data.id}`
    child.style.top = data.x
    child.style.left = data.y
    
    parent.appendChild(child)
}

function MoveTheCursor(data) {
    console.log(data)
    console.log(`move the cursor to (${data.x},${data.y})`)
    const cursor = document.getElementById(`cursor${data.id}`)
    cursor.style.left = `${data.x}px`
    cursor.style.top = `${data.y}px`
}

mouseMoveLastTime = 0
document.addEventListener("mousemove",(e)=>{
    const mouseMoveNow = Date.now();
    if (mouseMoveNow - mouseMoveLastTime < 50){
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

function updateCursorList(data){
    for (const item of CursorList){
        if (item.id != data.id) continue;
        item.x = data.x;
        item.y = data.y;
        MoveTheCursor(data);
        return
    }
    CursorList.push(data)
    CreateDOMCursor(data)
}

function reflashCursorList(dataList){
    // todo reflash
}

function deleteCursorList(id){
    // todo delete
}