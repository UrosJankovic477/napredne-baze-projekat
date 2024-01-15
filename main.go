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
	http.HandleFunc("/removeInterest", removeInterestHandler)
	http.HandleFunc("/createForum", createForumHandler)
	http.HandleFunc("/addPost", addPostHandler)
	http.HandleFunc("/getPost", getPostHandler)
	http.HandleFunc("/addComment", addCommentHandler)
	http.HandleFunc("/recommendForums", recommendForumHandler)
	http.HandleFunc("/recommendFriends", recommendFriendHandler)
	http.HandleFunc("/getPostsFromForum", getPostsFromForumHandler)
	http.HandleFunc("/getCommentsFromPost", getCommentsFromPostHandler)
	http.HandleFunc("/getPosts", getPostsHandler)

	http.ListenAndServe("localhost:8080", nil)
}
