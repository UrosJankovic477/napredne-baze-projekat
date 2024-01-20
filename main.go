package main

import (
	"napredne_baze_podataka/internals"
	"net/http"

	"github.com/philippseith/signalr"
)

func main() {
	internals.Hub = internals.ChatHub{}
	internals.Initialize()

	router := http.NewServeMux()

	router.HandleFunc("/api/register", registerHandler)
	router.HandleFunc("/api/login", loginHandler)
	router.HandleFunc("/api/logout", logoutHandler)
	router.HandleFunc("/api/friendRequest", friendRequestHandler)
	router.HandleFunc("/api/acceptRequest", acceptRequestHandler)
	router.HandleFunc("/api/declineRequest", declineRequestHandler)
	router.HandleFunc("/api/unfriend", unfriendHandler)
	router.HandleFunc("/api/addInterest", addInterestHandler)
	router.HandleFunc("/api/removeInterest", removeInterestHandler)
	router.HandleFunc("/api/createForum", createForumHandler)
	router.HandleFunc("/api/addPost", addPostHandler)
	router.HandleFunc("/api/getPost", getPostHandler)
	router.HandleFunc("/api/addComment", addCommentHandler)
	router.HandleFunc("/api/recommendForums", recommendForumHandler)
	router.HandleFunc("/api/recommendFriends", recommendFriendHandler)
	router.HandleFunc("/api/getPostsFromForum", getPostsFromForumHandler)
	router.HandleFunc("/api/getCommentsFromPost", getCommentsFromPostHandler)
	router.HandleFunc("/api/getPosts", getPostsHandler)
	router.HandleFunc("/api/makeChatRoom", makeChatRoomHandler)
	router.HandleFunc("/api/getUsersChatrooms", getUsersChatroomsHandler)
	router.HandleFunc("/api/getFriends", getFriendsHandler)
	router.HandleFunc("/api/getFriendRequests", getFriendRequestsHandler)
	router.HandleFunc("/api/getInterests", getInterestsHandler)
	router.HandleFunc("/api/searchForums", searchForumsHandler)
	router.HandleFunc("/api/searchPosts", searchPostsHandler)
	router.HandleFunc("/api/searchUsers", searchUsersHandler)
	router.HandleFunc("/api/uploadImage", uploadImageHandler)
	router.HandleFunc("/api/getImage/", getImageHandler)

	router.Handle("/", http.FileServer(http.Dir("wwwroot")))

	internals.Server.MapHTTP(signalr.WithHTTPServeMux(router), "/api/chat")

	http.ListenAndServe("localhost:8080", router)
}
