let intrests = []

let token = document.cookie.replace("token=", "");

let forumname =  document.querySelector("input#forumname")
let intreststxt = document.querySelector("input#intreststxt")
let intrestsul = document.querySelector("ul#intrestsul")
let addintrest = document.querySelector("button#addintrest")
let createbutton = document.querySelector("button#createbutton")

addintrest.onclick = () => {
    intrests.push(intreststxt.value);
    let li = document.createElement("li")
    let interest_txt = intreststxt.value;
    li.innerHTML = interest_txt;
    li.onmouseup = (e) => {
        if(e.button !== 2) {return}
        e.stopPropagation();
        e.preventDefault();
        let txt = e.target.value
        let tmp = intrests.filter(i => i !== txt)
        intrests = tmp
        e.target.remove()

    }
    createbutton.onclick = () => {
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        
        var raw = JSON.stringify({
          "UserToken": token,
          "Name": forumname.value,
          "Interests": intrests
        });
        
        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: raw,
          redirect: 'follow'
        };
        
        fetch("/api/createForum", requestOptions)
          .then(response => response.text())
          .then(result => console.log(result))
          .catch(error => console.log('error', error));
    }
    intreststxt.value = ''
    intrestsul.appendChild(li)


}