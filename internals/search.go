package internals

import (
	"net/http"
	"strings"
)

func SearchForums(search_querry string) ([]string, int, error) {
	result, err := doQuery("MATCH (forum:Forum) WHERE TOLOWER(forum.Name) CONTAINS $search_querry RETURN forum.Name AS forum_name",
		map[string]any{
			"search_querry": strings.ToLower(search_querry),
		})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	forums := make([]string, 0)
	for _, record := range result.Records {
		forum_name, _ := record.Get("forum_name")
		forums = append(forums, forum_name.(string))
	}
	return forums, http.StatusOK, nil
}

func SearchPosts(search_querry string) ([]Post, int, error) {
	result, err := doQuery("MATCH (post:Post) WHERE TOLOWER(post.Title) CONTAINS $search_querry RETURN post.Title AS post_title, post.UUID as post_uuid, post.PostedOn AS posted_on",
		map[string]any{
			"search_querry": strings.ToLower(search_querry),
		})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	posts := make([]Post, 0)
	for _, record := range result.Records {
		post_title, _ := record.Get("post_title")
		post_uuid, _ := record.Get("post_uuid")
		posted_on, _ := record.Get("posted_on")
		posts = append(posts, Post{Title: post_title.(string), UUID: post_uuid.(string), PostedOn: posted_on.(int64)})
	}
	return posts, http.StatusOK, nil
}

func SearchUsers(search_querry string) ([]string, int, error) {
	result, err := doQuery("MATCH (usr:AccountCredentials) WHERE TOLOWER(usr.Username) CONTAINS $search_querry RETURN usr.Username AS username",
		map[string]any{
			"search_querry": strings.ToLower(search_querry),
		})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	users := make([]string, 0)
	for _, record := range result.Records {
		username, _ := record.Get("username")
		users = append(users, username.(string))
	}
	return users, http.StatusOK, nil
}
