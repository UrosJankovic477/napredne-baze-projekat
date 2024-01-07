package main

import (
	"encoding/json"
	"io"
	"log"
	"napredne_baze_podataka/internals"
	"net/http"
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

	forums, status, err := internals.RecommendForums(deserialized["UserToken"])
	if err != nil {
		writer.WriteHeader(status)
		log.Println(err)
	}
	to_json, err := json.Marshal(forums)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	writer.Write(to_json)
}
