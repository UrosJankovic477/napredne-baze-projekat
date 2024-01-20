package main

import (
	"encoding/json"
	"io"
	"log"
	"napredne_baze_podataka/internals"
	"net/http"
	"strings"
)

type login_struct struct {
	Username string
	Password string
}

type friend_request_struct struct {
	UserToken  string
	FriendName string
}

type forum_creation struct {
	UserToken string
	Name      string
	Interests []string
}

func registerHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	acc := internals.AccountCredentials{}
	err = json.Unmarshal(body, &acc)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = internals.CreateAccount(acc.Username, acc.PasswordHash)
	if err != nil {
		writer.WriteHeader(http.StatusConflict)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func loginHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	deserialized := login_struct{}
	json.Unmarshal(body, &deserialized)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := internals.LoginUser(deserialized.Username, deserialized.Password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	token_json, err := json.Marshal(&token)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	writer.Write(token_json)
}

func logoutHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "DELETE" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	token := map[string]string{
		"Token": "",
	}
	json.Unmarshal(body, &token)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	status, err := internals.DeleteToken(token["Token"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func friendRequestHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken":  "",
		"FriendName": "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.SendFriendRequest(deserialized["UserToken"], deserialized["FriendName"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func acceptRequestHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken":  "",
		"FriendName": "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.AcceptRequest(deserialized["UserToken"], deserialized["FriendName"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func declineRequestHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "DELETE" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken":  "",
		"FriendName": "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.DeclineRequest(deserialized["UserToken"], deserialized["FriendName"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func unfriendHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken":  "",
		"FriendName": "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.Unfriend(deserialized["UserToken"], deserialized["FriendName"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func addInterestHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken": "",
		"Interest":  "",
		"Category":  "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.AddInterest(deserialized["UserToken"], deserialized["Category"], deserialized["Interest"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func removeInterestHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"Token":    "",
		"Interest": "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.RemoveInterest(deserialized["Token"], deserialized["Interest"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func createForumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := forum_creation{}
	json.Unmarshal(body, &deserialized)

	status, err := internals.CreateForum(deserialized.UserToken, deserialized.Name, deserialized.Interests)
	if err != nil {
		writer.WriteHeader(status)
		if status == http.StatusConflict {
			writer.Write([]byte(err.Error()))
		}
		log.Println(err)
	}
}

func addPostHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken": "",
		"ForumName": "",
		"Title":     "",
		"Body":      "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.AddPost(deserialized["UserToken"], deserialized["ForumName"], deserialized["Title"], deserialized["Body"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func addCommentHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken": "",
		"PostUUID":  "",
		"Body":      "",
	}
	json.Unmarshal(body, &deserialized)

	status, err := internals.CommentPost(deserialized["UserToken"], deserialized["PostUUID"], deserialized["Body"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
}

func recommendForumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken": "",
	}
	json.Unmarshal(body, &deserialized)

	recommendations, status, err := internals.RecommendForums(deserialized["UserToken"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(recommendations)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}

func recommendFriendHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UserToken": "",
	}
	json.Unmarshal(body, &deserialized)

	recommendations, status, err := internals.RecommendFriends(deserialized["UserToken"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(recommendations)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}

func getPostHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]string{
		"UUID": "",
	}
	json.Unmarshal(body, &deserialized)

	post, status, err := internals.GetPost(deserialized["UUID"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(post)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}

func getPostsFromForumHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"ForumName": "",
		"Limit":     0,
		"Offset":    0,
	}
	json.Unmarshal(body, &deserialized)

	posts, status, err := internals.GetMultiplePostsFromForum(
		deserialized["ForumName"].(string),
		int(deserialized["Limit"].(float64)),
		int(deserialized["Offset"].(float64)))
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(posts)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}

func getCommentsFromPostHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"ForumName": "",
		"Limit":     0,
		"Offset":    0,
	}
	json.Unmarshal(body, &deserialized)

	comments, status, err := internals.GetCommentsFromPost(
		deserialized["PostUUID"].(string),
		int(deserialized["Limit"].(float64)),
		int(deserialized["Offset"].(float64)))
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(comments)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}

func getPostsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"Token":  "",
		"Limit":  0,
		"Offset": 0,
	}
	json.Unmarshal(body, &deserialized)

	posts, status, err := internals.GetPosts(
		deserialized["Token"].(string),
		int(deserialized["Limit"].(float64)),
		int(deserialized["Offset"].(float64)))
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(posts)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}

type ChatRoomReqBody struct {
	Token string
	Name  string
	Users []string
}

func makeChatRoomHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := ChatRoomReqBody{}
	json.Unmarshal(body, &deserialized)

	status, err := internals.MakeChatroomNode(
		deserialized.Token,
		deserialized.Name,
		deserialized.Users)

	writer.WriteHeader(status)
}

func joinChatRoomHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"User": "",
		"UUID": "",
	}
	json.Unmarshal(body, &deserialized)
	status, err := internals.JoinChatroom(deserialized["User"].(string), deserialized["UUID"].(string))
	writer.WriteHeader(status)
}

func getUsersChatroomsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"Token": "",
	}
	json.Unmarshal(body, &deserialized)
	chatrooms, status, err := internals.GetUsersChatrooms(deserialized["Token"].(string))

	to_json, err := json.Marshal(chatrooms)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func getFriendsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"Token": "",
	}
	json.Unmarshal(body, &deserialized)
	friends, status, err := internals.GetFriends(deserialized["Token"].(string))

	to_json, err := json.Marshal(friends)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func getFriendRequestsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"Token": "",
	}
	json.Unmarshal(body, &deserialized)
	friendRequests, status, err := internals.GetFriendRequests(deserialized["Token"].(string))

	to_json, err := json.Marshal(friendRequests)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func getInterestsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"Token": "",
	}
	json.Unmarshal(body, &deserialized)
	interests, status, err := internals.GetInterests(deserialized["Token"].(string))

	to_json, err := json.Marshal(interests)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func searchForumsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"SearchQuery": "",
	}
	json.Unmarshal(body, &deserialized)
	interests, status, err := internals.SearchForums(deserialized["SearchQuery"].(string))

	to_json, err := json.Marshal(interests)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func searchUsersHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"SearchQuery": "",
	}
	json.Unmarshal(body, &deserialized)
	interests, status, err := internals.SearchUsers(deserialized["SearchQuery"].(string))

	to_json, err := json.Marshal(interests)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func searchPostsHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	deserialized := map[string]any{
		"SearchQuery": "",
	}
	json.Unmarshal(body, &deserialized)
	interests, status, err := internals.SearchPosts(deserialized["SearchQuery"].(string))

	to_json, err := json.Marshal(interests)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func uploadImageHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(reqptr.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	interests, status, err := internals.UploadImage(body)

	to_json, err := json.Marshal(interests)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Write(to_json)
}

func getImageHandler(writer http.ResponseWriter, reqptr *http.Request) {
	if reqptr.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hash := strings.TrimPrefix(reqptr.URL.Path, "/api/getImage/")

	image, status, err := internals.GetImage(hash)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	writer.WriteHeader(status)
	writer.Header().Add("Content-Type", http.DetectContentType(image))
	writer.Write(image)
}
