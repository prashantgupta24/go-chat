window.onload = function () {
    let usernameField = document.getElementById("username");
    let username = prompt("Please enter your name", "Harry Potter");
    if (!username) {
        return false
    }
    usernameField.value = username

    let conn;
    let msg = document.getElementById("msg");
    let log = document.getElementById("log");

    function appendLog(item) {
        let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        let messageToSend = {
            type: "msg",
            sender: usernameField.value,
            message: msg.value,
            date: Date.now()
        };
        conn.send(JSON.stringify(messageToSend))
        console.log("Sending msg : " + JSON.stringify(messageToSend))
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");

        conn.onopen = function (event) {
            let messageToSend = {
                type: "reg",
                sender: usernameField.value,
                message: "",
                date: Date.now()
            };
            conn.send(JSON.stringify(messageToSend))
            console.log("Sending reg : " + JSON.stringify(messageToSend))
            msg.focus()
        };

        conn.onclose = function (evt) {
            let item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            let messageReceivedJSON = JSON.parse(evt.data);
            if (messageReceivedJSON.type == "reg") {
                console.log("received reg->" + evt.data)
                let item = document.createElement("div");
                item.innerText = "**** " + messageReceivedJSON.sender + " **** joined. Total participants : " + messageReceivedJSON.message;
                appendLog(item);
            } else if (messageReceivedJSON.type == "unreg") {
                console.log("received unreg->" + evt.data)
                let item = document.createElement("div");
                item.innerText = "**** " + messageReceivedJSON.sender + " **** left. Total participants : " + messageReceivedJSON.message;
                appendLog(item);
            } else {
                console.log("received msg->" + evt.data)
                let item = document.createElement("div");
                item.innerText = "< " + messageReceivedJSON.sender + " > : " + messageReceivedJSON.message;
                appendLog(item);

            }
        };
    } else {
        let item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};