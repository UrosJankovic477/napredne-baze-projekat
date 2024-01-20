let results = document.getElementById("results");
let token = document.cookie.replace("token=", "");
let post_uuid = window.location.hash.substring(1);
let post_div = document.getElementById("post_title");
let post_body_div = document.getElementById("post_body");
let comments_list = document.getElementById("comments");
let comment_txtbox = document.getElementById("comment_txtbox");
let add_comment_btn = document.getElementById("add_comment_btn");

var myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json")
var raw = JSON.stringify({
  "UUID": post_uuid
})
var requestOptions = {
  method: 'POST',
  headers: myHeaders,
  body: raw,
  redirect: 'follow'
}
fetch("/api/getPost", requestOptions)
  .then(response => response.json())
  .then(result => {
    post_div.innerHTML = result.Title;
    post_body_div.innerHTML = result.Body;
  })
  .catch(error => console.log('error', error));

var myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");

var raw = JSON.stringify({
  "PostUUID": post_uuid,
  "Limit": 65535.0,
  "Offset": 0.0
});

var requestOptions = {
  method: 'POST',
  headers: myHeaders,
  body: raw,
  redirect: 'follow'
};

fetch("/api/getCommentsFromPost", requestOptions)
  .then(response => response.json())
  .then(result => result.forEach(comment => {
    let li = document.createElement("li");

    let datum = new Date(comment.PostedOn * 1000).toLocaleDateString("sr-RS");
    let vreme = new Date(comment.PostedOn * 1000).toLocaleTimeString("sr-RS");

    li.innerHTML = `Posted by: ${comment.Author} at ${datum} : ${vreme}`;
    let br = document.createElement("br");
    li.appendChild(br);
    let comment_body_div = document.createElement("div");
    comment_body_div.innerHTML = comment.Body;
    li.appendChild(comment_body_div);
    comments_list.appendChild(li);
  }))
  .catch(error => console.log('error', error));

add_comment_btn.onclick = () => {
  let comment = comment_txtbox.value;
  var myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  var raw = JSON.stringify({
    "UserToken": token,
    "PostUUID": post_uuid,
    "Body": comment
  });
  var requestOptions = {
    method: 'POST',
    headers: myHeaders,
    body: raw,
    redirect: 'follow'
  };
  fetch("/api/addComment", requestOptions)
    .then(response => response.json())
    .then(result => console.log(result))
    .catch(error => console.log('error', error));
}
