"use strict";

let token = document.cookie.replace("token=", "");
let connection = new signalR.HubConnectionBuilder().withUrl("api/chat").build();
let myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");

let chatrooms = undefined;
let currentChatroomUUID

var raw = JSON.stringify({
  "Token": token
});

var requestOptions = {
  method: 'POST',
  headers: myHeaders,
  body: raw,
  redirect: 'follow'
};


window.onload = () => {
  fetch("/api/getUsersChatrooms", requestOptions)
  .then(response => response.text())
    .then(result => {
      chatrooms = JSON.parse(result);
      chatrooms.forEach(chatroom => {


        let li = document.createElement("input");
        li.type = "radio";
        li.value = chatroom.UUID;
        li.name = "chatroomRadio"
        li.onchange = (event) => {
          currentChatroomUUID = event.target.value;

          document.getElementById("messagesList").innerHTML = "";

          connection.invoke("JoinChat", currentChatroomUUID).catch(function (err) {
            return console.error(err.toString());
          });

        }
        let label = document.createElement("label");

        let name = undefined;
        if(chatroom.Name === "") {
          chatroom.Users.forEach(user => {
            name += user + " ";
            
          })
        }
        else {
          name = chatroom.Name;
        }

        label.innerHTML = name;
        chatroomList.appendChild(label);

        document.getElementById("chatroomList").appendChild(li);
      });
  
    })
    .catch(error => console.log('error', error));
  
}




//Disable the send button until connection is established.
document.getElementById("sendButton").disabled = true;

connection.on("ReceiveMessage", function (user, message, time) {
    var li = document.createElement("li");
    document.getElementById("messagesList").appendChild(li);

    let datum = new Date(time * 1000).toLocaleDateString("sr-RS");
    let vreme = new Date(time * 1000).toLocaleTimeString("sr-RS");

    li.textContent = `${user} says ${message} at ${datum} : ${vreme}`;
});

connection.on("ReceiveMessageList", function (msgs) {
  msgs.forEach(msg => {
    var li = document.createElement("li");
    document.getElementById("messagesList").appendChild(li);

    let datum = new Date(msg.Time * 1000).toLocaleDateString("sr-RS");
    let vreme = new Date(msg.Time * 1000).toLocaleTimeString("sr-RS");
    
    li.textContent = `${msg.Username} says ${msg.Content} at ${datum} : ${vreme}`;
  })

});

connection.start().then(function () {
    document.getElementById("sendButton").disabled = false;
}).catch(function (err) {
    return console.error(err.toString());
});

let users = {}


document.getElementById("addUser").addEventListener("click", function (event) {
  var username = document.getElementById("username").value;
  
  users[username] = true;
  var username = document.getElementById("username").value = "";

  event.preventDefault();
});


document.getElementById("makeChatRoom").addEventListener("click", function (event) {
  var chatroom_name = document.getElementById("groupName").value;
  if (Object.keys(users).length === 0) {
    console.error("There must be more than one user in a chat");
  }
  else {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    var raw = JSON.stringify({
        "Token": token,
        "Name": chatroom_name,
        "Users": Object.keys(users)
      });
  
    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };
    fetch("/api/makeChatRoom", requestOptions)
      .then(response => response.text())
      .then(result => {
        console.log(result);
        document.getElementById("groupName").value = "";
      })
      .catch(error => console.log('error', error));
    }



  event.preventDefault();
});


document.getElementById("sendButton").addEventListener("click", function (event) {
    var message = document.getElementById("messageInput").value;
    connection.invoke("SendMessage", token, message, currentChatroomUUID).catch(function (err) {
        return console.error(err.toString());
    });
    event.preventDefault();
});


