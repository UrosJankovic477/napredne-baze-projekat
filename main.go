package main

import (
	"napredne_baze_podataka/internals"
	"net/http"
)

func main() {
	internals.Initialize()

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/friendRequest", friendRequestHandler)
	http.HandleFunc("/acceptRequest", acceptRequestHandler)
	http.HandleFunc("/declineRequest", declineRequestHandler)
	http.HandleFunc("/unfriend", unfriendHandler)
	http.HandleFunc("/addInterest", addInterestHandler)
	http.ListenAndServe("localhost:8080", nil)

}
