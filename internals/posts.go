package internals

import (
	"net/http"
	"sort"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type Post struct {
	Author       string
	Title        string
	Body         string
	UUID         string
	PostedOn     int64
	CommentCount int64
}

type Comment struct {
	Author   string
	Body     string
	PostedOn int64
}

func AddToRedis(UUID string, body string, isPost bool) error {
	hash_name := ""

	if isPost {
		hash_name = "post"
	} else {
		hash_name = "comment"
	}

	err := rdb.HSet(ctx, hash_name, UUID, body).Err()
	return err
}

func GetFromRedis(UUID string, isPost bool) (string, error) {
	hash_name := ""

	if isPost {
		hash_name = "post"
	} else {
		hash_name = "comment"
	}

	val, err := rdb.HGet(ctx, hash_name, UUID).Result()

	return val, err
}

func GetMultiplePostsFromForum(forum_name string, count int, offset int) ([]Post, int, error) {
	result, err := doQuery("MATCH (forum:Forum) WHERE forum.Name = $forum_name "+
		"WITH forum MATCH (posts:Post) WHERE (forum)-[:HAS_POST]->(posts) "+
		"WITH posts ORDER BY posts.PostedOn "+
		"OPTIONAL MATCH (posts)-[r:HAS_COMMENT]->(comment) "+
		"WITH posts, COUNT(r) AS comment_count ORDER BY comment_count "+
		"MATCH (posts)-[posted_by:POSTED_BY]-(author:AccountCredentials) "+
		"RETURN posts, author.Username AS author, comment_count SKIP $offset LIMIT $count ",
		map[string]any{
			"forum_name": forum_name,
			"count":      count,
			"offset":     offset,
		})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	posts := make([]Post, 0)
	for _, record := range result.Records {
		post_node, _ := record.Get("posts")
		author, _ := record.Get("author")
		comment_count, _ := record.Get("comment_count")
		title := post_node.(dbtype.Node).Props["Title"]
		UUID := post_node.(dbtype.Node).Props["UUID"]
		PostedOn := post_node.(dbtype.Node).Props["PostedOn"]
		//body, err := GetFromRedis(UUID.(string), true)
		//if err != nil {
		//	return nil, http.StatusNotFound, err
		//}
		posts = append(posts, Post{
			Title:        title.(string),
			Author:       author.(string),
			UUID:         UUID.(string),
			PostedOn:     PostedOn.(int64),
			CommentCount: comment_count.(int64)})
	}
	return posts, http.StatusOK, nil
}

func GetCommentsFromPost(UUID string, count int, offset int) ([]Comment, int, error) {
	result, err := doQuery("MATCH (post:Post) WHERE post.UUID = $UUID "+
		"WITH post MATCH (comments:Comment) WHERE (post)-[:HAS_COMMENT]->(comments) "+
		"WITH comments ORDER BY comments.PostedOn "+
		"MATCH (comments)-[posted_by:POSTED_BY]-(author:AccountCredentials) "+
		"RETURN comments, author.Username AS author SKIP $offset LIMIT $count ",
		map[string]any{
			"UUID":   UUID,
			"count":  count,
			"offset": offset,
		})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	comments := make([]Comment, 0)
	for _, record := range result.Records {
		comment_node, _ := record.Get("comments")
		author, _ := record.Get("author")
		UUID := comment_node.(dbtype.Node).Props["UUID"]
		PostedOn := comment_node.(dbtype.Node).Props["PostedOn"]
		body, err := GetFromRedis(UUID.(string), false)
		if err != nil {
			return nil, http.StatusNotFound, err
		}
		comments = append(comments, Comment{
			Author:   author.(string),
			PostedOn: PostedOn.(int64),
			Body:     body})
	}

	return comments, http.StatusOK, nil

}

func GetPost(UUID string) (Post, int, error) {
	result, err := doQuery("MATCH (post) WHERE post.UUID = $UUID "+
		"WITH post "+
		"OPTIONAL MATCH (post)-[r:HAS_COMMENT]->(comment) "+
		"WITH post, COUNT(r) AS comment_count ORDER BY comment_count "+
		"MATCH (post)-[posted_by:POSTED_BY]-(author:AccountCredentials) "+
		"RETURN post, author.Username AS author, comment_count",
		map[string]any{
			"UUID": UUID,
		})
	if err != nil {
		return Post{}, http.StatusInternalServerError, err
	}

	body, err := GetFromRedis(UUID, true)
	if err != nil {
		return Post{}, http.StatusInternalServerError, err
	}

	record := result.Records[0]
	post_node, _ := record.Get("post")
	author, _ := record.Get("author")
	comment_count, _ := record.Get("comment_count")
	title := post_node.(dbtype.Node).Props["Title"]
	PostedOn := post_node.(dbtype.Node).Props["PostedOn"]

	post := Post{
		Title:        title.(string),
		Author:       author.(string),
		Body:         body,
		PostedOn:     PostedOn.(int64),
		UUID:         UUID,
		CommentCount: comment_count.(int64)}

	return post, http.StatusOK, nil
}

func GetPosts(token string, count int, offset int) ([]Post, int, error) {
	_, _, err := GetUserFromToken(token)
	if err != nil {
		// preporuci nepersonalizovane postove
		return nil, http.StatusNotImplemented, err
	}

	forums, status, err := RecommendForums(token)
	if err != nil {
		return nil, status, err
	}

	posts := make([]Post, 0)
	var result *neo4j.EagerResult
	for _, forum := range forums {
		result, err = doQuery("MATCH (forum:Forum) WHERE forum.Name = $forum_name "+
			"WITH forum MATCH (posts:Post) WHERE (forum)-[:HAS_POST]->(posts) "+
			"WITH posts ORDER BY posts.PostedOn "+
			"OPTIONAL MATCH (posts)-[r:HAS_COMMENT]->(comment) "+
			"WITH posts, COUNT(r) AS comment_count ORDER BY comment_count "+
			"MATCH (posts)-[posted_by:POSTED_BY]-(author:AccountCredentials) "+
			"RETURN posts, author.Username AS author, comment_count SKIP $offset LIMIT $count ",
			map[string]any{
				"forum_name": forum.Name,
				"count":      count,
				"offset":     offset,
			})
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		for _, record := range result.Records {
			post_node, _ := record.Get("posts")
			author, _ := record.Get("author")
			comment_count, _ := record.Get("comment_count")
			title := post_node.(dbtype.Node).Props["Title"]
			UUID := post_node.(dbtype.Node).Props["UUID"]
			PostedOn := post_node.(dbtype.Node).Props["PostedOn"]
			//body, err := GetFromRedis(UUID.(string), true)
			//if err != nil {
			//	return nil, http.StatusNotFound, err
			//}
			posts = append(posts, Post{
				Title:        title.(string),
				Author:       author.(string),
				UUID:         UUID.(string),
				PostedOn:     PostedOn.(int64),
				CommentCount: comment_count.(int64)})
		}

	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CommentCount > posts[j].CommentCount
	})

	return posts, http.StatusOK, nil
}
