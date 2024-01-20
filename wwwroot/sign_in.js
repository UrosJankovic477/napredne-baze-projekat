let username = document.getElementById("usernameInput")
let password = document.getElementById("passwordInput")
let login = document.getElementById("logInButton");

login.onclick = () => {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    var raw = JSON.stringify({
      "Username": username.value,
      "Password": password.value
    });

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    fetch("/api/login", requestOptions)
      .then(response => response.text())
      .then(result => {
        document.cookie = `token=${JSON.parse(result)}`;
        window.location.href = "/mainpage.html"
      })
      .catch(error => console.log('error', error));
}

var signin = document.getElementById("signInButton");
signin.onclick = () => {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");


    var raw = JSON.stringify({
      "Username": username.value,
      "PasswordHash": password.value
    });

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    fetch("/api/register", requestOptions)
      .then(response => response.text())
      .then(result => log(result))
      .catch(error => console.log('error', error));
}
