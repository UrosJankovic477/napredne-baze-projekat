"use strict";

let token = document.cookie.replace("token=", "");
let connection = new signalR.HubConnectionBuilder().withUrl("api/chat").build();
let myHeaders = new Headers();
let upload_img_btn = document.getElementById("upload_img_btn");
let current_img_hash = undefined;
let img_upload = document.createElement("input");
myHeaders.append("Content-Type", "application/json");
img_upload.onsubmit = (e) => {
  e.preventDefault()
}
upload_img_btn.onclick = () => {

  
  img_upload.type = "file";
  img_upload.click()

 
}

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
  let idx = message.lastIndexOf(" EMBED-IMAGE:")
  if (idx != -1) {
    let embed_string = message.substring(idx)
    if (embed_string != undefined) {
    let img = document.createElement("img");
    let replaced = embed_string.replace(" EMBED-IMAGE:", "")
    img.src = `/api/getImage/${replaced}`
    document.getElementById("messagesList").appendChild(img)
    message = message.replace(embed_string, "")
    }
  }
  

  var li = document.createElement("li");
  document.getElementById("messagesList").appendChild(li);
  let datum = new Date(time * 1000).toLocaleDateString("sr-RS");
  let vreme = new Date(time * 1000).toLocaleTimeString("sr-RS");
  li.textContent = `${user} says ${message} at ${datum} : ${vreme}`;
});

connection.on("ReceiveMessageList", function (msgs) {
  msgs.forEach(msg => {
    
    let idx = msg.Content.lastIndexOf(" EMBED-IMAGE:")
    if (idx != -1) {
      let embed_string = msg.Content.substring(idx)
      if (embed_string != undefined) {
      let img = document.createElement("img");
      let replaced = embed_string.replace(" EMBED-IMAGE:", "")
      img.src = `/api/getImage/${replaced}`
      document.getElementById("messagesList").appendChild(img)
      msg.Content = msg.Content.replace(embed_string, "")
      }
    }
  

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
    if (img_upload.value !== "") {
    //  img_upload.remove()
    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      redirect: 'follow'
    };
  
    requestOptions.body = img_upload.files[0]
    
    fetch("/api/uploadImage", requestOptions)
    .then(response => response.json())
    .then(result => {
      message += ` EMBED-IMAGE:${result}`
      img_upload.value = ""
      connection.invoke("SendMessage", token, message, currentChatroomUUID).catch(function (err) {
        return console.error(err.toString());
    });
    })
    .catch(error => console.error(error))


    }
    else {
      connection.invoke("SendMessage", token, message, currentChatroomUUID).catch(function (err) {
        return console.error(err.toString());
    });
    }

    
    
    
    event.preventDefault();
});


