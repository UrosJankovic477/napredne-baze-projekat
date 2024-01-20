let token = document.cookie.replace("token=", "");

let recommend_friends_btn = document.getElementById("recommend_friends_btn");
let recommend_forums_btn = document.getElementById("recommend_forums_btn");
let results =  document.getElementById("results");

var myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");


var raw = JSON.stringify({
  "Token": token,
  "Limit": 65535.0,
  "Offset": 0.0
});

var requestOptions = {
  method: 'POST',
  headers: myHeaders,
  body: raw,
  redirect: 'follow'
};

fetch("/api/getPosts", requestOptions)
  .then(response => response.json())
  .then(result => result.forEach(post => {
    let li = document.createElement("li");

    let datum = new Date(post.PostedOn * 1000).toLocaleDateString("sr-RS");
    let vreme = new Date(post.PostedOn * 1000).toLocaleTimeString("sr-RS");

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
  .catch(error => console.log('error', error));


recommend_friends_btn.onclick = () => {
    results.innerHTML = "";
    let rec_friends;
    
    
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
      "UserToken": token
    });

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    fetch("/api/recommendFriends", requestOptions)
      .then(response => response.json())
      .then(result => result.forEach(rec_friend => {

        let li = document.createElement("li");
        let rec_friend_div = document.createElement("div");
        li.appendChild(rec_friend_div);
        rec_friend_div.innerHTML = rec_friend.Username;
        let friend_req_btn = document.createElement("button");
        friend_req_btn.innerHTML = "Send Friend Request";
        friend_req_btn.onclick = () => {
    
            var myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");

            var raw = JSON.stringify({
              "UserToken": token,
              "FriendName": rec_friend.Username
            });

            var requestOptions = {
              method: 'POST',
              headers: myHeaders,
              body: raw,
              redirect: 'follow'
            };

            fetch("/api/friendRequest", requestOptions)
              .then(response => response.text())
              .then(result => console.log(result))
              .catch(error => console.log('error', error));
        }
        rec_friend_div.appendChild(friend_req_btn);
        rec_friend_div.appendChild(document.createElement("br"));
        let rec_reason = "";
        let interests = Object.keys(rec_friend.Interests);
        if (interests.length > 0) {
            rec_reason += "Because you both like"
            interests.forEach(interest => {
                rec_reason += " " + interest;
            })
        }
        let friends = Object.keys(rec_friend.Friends);
        if (friends.length > 0) {
            if (interests.length > 0) {
                rec_reason += " and ";
            }
            rec_reason += "because you're both friends with"
            friends.forEach(friend => {
                rec_reason += " " + friend;
            })
        }
        let reason_div = document.createElement("div");
        reason_div.innerHTML = rec_reason;
        rec_friend_div.appendChild(reason_div);
        results.appendChild(rec_friend_div);
      }))
      .catch(error => console.log('error', error));

}
recommend_forums_btn.onclick = () => {
    results.innerHTML = "";
    let rec_forums;

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
      "UserToken": token
    });

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };
    fetch("/api/recommendForums", requestOptions)
      .then(response => response.json())
      .then(result => result.forEach(rec_forum => {
        let li = document.createElement("li");
        let rec_forum_a = document.createElement("a");
        li.appendChild(rec_forum_a);
        rec_forum_a.innerHTML = rec_forum.Name;
        rec_forum_a.href = `/forum.html#${rec_forum.Name}`;

            
        rec_forum_a.appendChild(document.createElement("br"));
        let rec_reason = "";
        let interests = Object.keys(rec_forum.Interests);
        if (interests.length > 0) {
            rec_reason += "Because you like"
            interests.forEach(interest => {
                rec_reason += " " + interest;
            })
        }
        let friends = Object.keys(rec_forum.Friends);
        if (friends.length > 0) {
            if (interests.length > 0) {
                rec_reason += " and ";
            }
            rec_reason += "because "
            friends.forEach(friend => {
                rec_reason += " " + friend;
            })
            rec_reason += " are also active here";
        }
        let reason_div = document.createElement("div");
        reason_div.innerHTML = rec_reason;
        rec_forum_a.appendChild(reason_div);
        results.appendChild(rec_forum_a);

      }))
      .catch(error => console.log('error', error));
}

let friendrequestdiv = document.querySelector("div#friendrequestsdiv")

 window.onload = () => {

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
      "Token": token
    });

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    fetch("/api/getFriendRequests", requestOptions)
      .then(response => response.json())
      .then(result => result.forEach(r => {

         let reqdiv =  document.createElement("div")
         let username = document.createElement("p")
         username.innerHTML = r
         reqdiv.appendChild(username)
         let accept_button = document.createElement("button")
         accept_button.innerHTML = "Accept"
         accept_button.onclick = () => {
          var myHeaders = new Headers();
          myHeaders.append("Content-Type", "application/json");

          var raw = JSON.stringify({
            "UserToken": token,
            "FriendName": r
          });

          var requestOptions = {
            method: 'POST',
            headers: myHeaders,
            body: raw,
            redirect: 'follow'
          };

          fetch("/api/acceptRequest", requestOptions)
            .then(response => response.text())
            .then(result => console.log(result))
            .catch(error => console.log('error', error));
                   }
         let decline_button = document.createElement("button")
         decline_button.innerHTML = "Decline"
         decline_button.onclick = () => {
          var myHeaders = new Headers();
          myHeaders.append("Content-Type", "application/json");
          
          var raw = JSON.stringify({
            "UserToken": token,
            "FriendName": r
          });
          
          var requestOptions = {
            method: 'POST',
            headers: myHeaders,
            body: raw,
            redirect: 'follow'
          };
          
          fetch("/api/declineRequest", requestOptions)
            .then(response => response.text())
            .then(result => console.log(result))
            .catch(error => console.log('error', error));
         }
         reqdiv.appendChild(username)
         reqdiv.appendChild(accept_button)
         reqdiv.appendChild(decline_button)
         friendrequestdiv.appendChild(reqdiv)
 
     }))
      .catch(error => console.log('error', error));
     



      let my_friends = document.getElementById("my_friends");

      fetch("/api/getFriends", requestOptions)
      .then(response => response.json())
      .then(result => result.forEach(r => {

         let friendDiv =  document.createElement("div")
         let username = document.createElement("p")
         username.innerHTML = r
         friendDiv.appendChild(username)
         let unfriend_button = document.createElement("button")
         unfriend_button.innerHTML = "Unfriend"
         unfriend_button.onclick = () => {
          var myHeaders = new Headers();
          myHeaders.append("Content-Type", "application/json");

          var raw = JSON.stringify({
            "UserToken": token,
            "FriendName": r
          });

          var requestOptions = {
            method: 'POST',
            headers: myHeaders,
            body: raw,
            redirect: 'follow'
          };

          fetch("/api/unfriend", requestOptions)
            .then(response => response.text())
            .then(result => console.log(result))
            .catch(error => console.log('error', error));
                   }
        
         friendDiv.appendChild(username)
         friendDiv.appendChild(unfriend_button)

         my_friends.appendChild(friendDiv)
 
     }))
      .catch(error => console.log('error', error));


      let interests_div = document.getElementById("interests");
      let add_interest_btn = document.getElementById("add_interest_btn");
      let interest_txtbox = document.getElementById("interest_txtbox")
      add_interest_btn.onclick = () => {
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        
        var raw = JSON.stringify({
          "UserToken": token,
          "Interest": interest_txtbox.value,
          "Category": ""
        });
        
        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: raw,
          redirect: 'follow'
        };
        
        fetch("/api/addInterest", requestOptions)
          .then(response => response.text())
          .then(result => console.log(result))
          .catch(error => console.log('error', error));
      }


      var myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/json");

      var raw = JSON.stringify({
        "Token": token
      });

      var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
      };

      fetch("/api/getInterests", requestOptions)
        .then(response => response.json())
        .then(result => result.forEach(interest => {
          let li = document.createElement("li");
          li.innerHTML = interest;
          let remove_interest_btn = document.createElement("button");
          remove_interest_btn.innerHTML = "Remove";
          remove_interest_btn.onclick = () => {
            var myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");
            
            var raw = JSON.stringify({
              "Token": token,
              "Interest": interest
            });
            
            var requestOptions = {
              method: 'POST',
              headers: myHeaders,
              body: raw,
              redirect: 'follow'
            };
            
            fetch("/api/removeInterest", requestOptions)
              .then(response => response.text())
              .then(result => console.log(result))
              .catch(error => console.log('error', error));
          }
          li.appendChild(remove_interest_btn);
          interests_div.appendChild(li);
        }))
        .catch(error => console.log('error', error));

    
 }
let opt = undefined;
let searchopts = document.getElementById("searchopts");
searchopts.onchange = (event) => {
  opt = event.target.value;
}
let searchbar = document.getElementById("searchbar");
let searchbtn = document.getElementById("searchbtn");
searchbtn.onclick = () => {
  let searchquery = searchbar.value;
  var myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  
  var raw = JSON.stringify({
    "SearchQuery": searchquery
  });
  
  var requestOptions = {
    method: 'POST',
    headers: myHeaders,
    body: raw,
    redirect: 'follow'
  };

  results.innerHTML = "";
  if (opt === "searchPosts") {
    fetch(`/api/${opt}`, requestOptions)
    .then(response => response.json())
    .then(result => result.forEach(res => {
      let li = document.createElement("li");
      let datum = new Date(res.PostedOn * 1000).toLocaleDateString("sr-RS");
      let vreme = new Date(res.PostedOn * 1000).toLocaleTimeString("sr-RS");

      li.innerHTML = `Posted by: ${res.Author}, at: ${datum} : ${vreme}`;
      let a = document.createElement("a");
      let br = document.createElement("br");
      li.appendChild(br);
      a.innerHTML = res.Title;
      a.href = `post.html#${res.UUID}`;
      li.appendChild(a);

      results.appendChild(li);
    }))
    .catch(error => console.log('error', error));
  }
  else if(opt === "searchForums") {
    fetch(`/api/${opt}`, requestOptions)
    .then(response => response.json())
    .then(result => result.forEach(res => {
      let li = document.createElement("li");

      let a = document.createElement("a");
      a.innerHTML = res;
      a.href = `forum.html#${res}`;
      li.appendChild(a);

      results.appendChild(li);
    }))
    .catch(error => console.log('error', error));
  }
  else if (opt === "searchUsers") {
    fetch(`/api/${opt}`, requestOptions)
    .then(response => response.json())
    .then(result => result.forEach(res => {
      let li = document.createElement("li");
      li.innerHTML = res;
      let btn = document.createElement("button");
      btn.innerHTML = "Send Friend Request";

      btn.onclick = () => {
    
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");

        var raw = JSON.stringify({
          "UserToken": token,
          "FriendName": res
        });

        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: raw,
          redirect: 'follow'
        };

        fetch("/api/friendRequest", requestOptions)
          .then(response => response.text())
          .then(result => console.log(result))
          .catch(error => console.log('error', error));
    }
    li.appendChild(btn);
    results.appendChild(li);
    }))
    .catch(error => console.log('error', error));
  }
}