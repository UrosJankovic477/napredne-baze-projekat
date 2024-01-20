let results = document.getElementById("results");
let token = document.cookie.replace("token=", "");
let forum_name = window.location.hash.substring(1);

window.onload = () => {
    let add_post_btn = document.getElementById("add_post_btn");
    add_post_btn.onclick = () => {
        window.location.href = `/make_post.html#${forum_name}`;
    }
    results.innerHTML = "";
    results.appendChild(add_post_btn)
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json")
    var raw = JSON.stringify({
      "ForumName": forum_name,
      "Limit": 65535.0,
      "Offset": 0.0 
    })
    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    }
    fetch("/api/getPostsFromForum", requestOptions)
      .then(response => response.json())
      .then(result => result.forEach(post => {
        let li = document.createElement("li");

        let datum = new Date(comment.PostedOn * 1000).toLocaleDateString("sr-RS");
        let vreme = new Date(comment.PostedOn * 1000).toLocaleTimeString("sr-RS");

        li.innerHTML =  `Posted by: ${post.Author}, at: ${datum} : ${vreme}`;
        let br = document.createElement("br");
        li.appendChild(br);
        let title_a = document.createElement("a");
        title_a.innerHTML = post.Title;
        title_a.href = `/post.html#${post.UUID}`;
        li.appendChild(title_a);
        comment_count_div = document.createElement("div");
        results.appendChild(li);
      }))
      .catch(error => console.log('error', error))
}