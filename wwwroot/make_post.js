let token = document.cookie.replace("token=", "");
let forum_name = window.location.hash.substring(1);


let post_btn = document.querySelector("button#post")

post_btn.onclick = () => {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    let post_body = document.querySelector("input#postbody").value
    let post_title = document.querySelector("input#postitle").value

    var raw = JSON.stringify({
      "UserToken": token,
      "ForumName": forum_name,
      "Title": post_title,
      "Body": post_body
    });

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    fetch("/api/addPost", requestOptions)
      .then(response => response.text())
      .then(result => console.log(result))
      .catch(error => console.log('error', error));
}