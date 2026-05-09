const Host = window.location.host
CursorList = []

var WSClient = new WebSocket(`ws://${Host}/api/cursor`)
clientID = ""

WSClient.onmessage = function (event) {
    // console.log("recieve: ",event.data)
    const data = JSON.parse(event.data)
    switch (data.status){
        case 1:
            updateCursorList(data.data);
            break;
        case 3:
            reflashCursorList(data.data);
            break;
        case 6:
            UpdateClientID(data.data);
            break;
        case 7:
            reflashCursorList(data.data);
            sendCursorCheck()
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
            x: 9999,
            y: 9999,
        }
    };
    sendMsgToWSClient(sendData);
}

function CreateDOMCursor(data){
    console.log(data)
    const parent = document.getElementById("cursor-list");
    const child = document.createElement("div");

    child.className = "cursor-item";
    child.id = `cursor${data.id}`
    // child.style.top = data.x
    // child.style.left = data.y
    
    parent.appendChild(child)
    MoveTheCursor(data)
}

function MoveTheCursor(data) {
    // console.log(`move the cursor to (${data.x},${data.y})`)
    const cursor = document.getElementById(`cursor${data.id}`)
    cursor.style.left = `${data.x}px`
    cursor.style.top = `${data.y}px`
}

function DeleteDOMCursor(data){
    console.log("delete: ", data)
    const child = document.getElementById(`cursor${data.id}`)

    child.remove()
}

var mouseMoveLastTime = 0
var cursorX = 0,cursorY = 0
document.addEventListener("mousemove",(e)=>{
    const mouseMoveNow = Date.now();
    if (mouseMoveNow - mouseMoveLastTime < 50){
        return;
    }
    mouseMoveLastTime = mouseMoveNow;
    cursorX = e.x;
    cursorY = e.y;
    // console.log(`mouse is in (${e.x},${e.y})`);
    
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
    // console.log(`already send msg: ${sendMsg}`)
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
    console.debug("reflash the cursor list")
    for (const data of dataList){
        existItem = false
        for(const item of CursorList){
            if (item.id != data.id) continue;
            existItem = true
        }
        
        if (existItem == false){
            CursorList.push(data)
            CreateDOMCursor(data)
        }
    }
    
    backList = []
    for (item of CursorList){
        existItem = false
        for(const data of dataList){
            if (item.id != data.id) continue;
            existItem = true
        }
        
        if (existItem ==  true) {
            backList.push(item)
        }else {
            DeleteDOMCursor(item)
        }
    }
    CursorList = backList
}

function deleteCursorList(id){
    // todo delete
}

function sendCursorCheck(){
    console.log("check point")
    const msgData = {
        status: 7,
        id: clientID,
        data: {
            id: clientID,
            x: cursorX,
            y: cursorY
        }
    }
    sendMsgToWSClient(msgData)
}