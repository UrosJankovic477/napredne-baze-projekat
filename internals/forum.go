package internals

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateForum(token string, name string, interests []string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	dict := make(map[string]any)
	dict["Id"] = user_node.ElementId
	dict["name"] = name
	query := ""
	for idx, interest := range interests {
		q := fmt.Sprintf("MERGE (interest_%d:Interest{Name:$%d}) ", idx, idx)
		query += q + fmt.Sprintf("MERGE (f)-[:HOSTS]->(interest_%d) ", idx)
		dict[fmt.Sprint(idx)] = interest
	}
	_, err = doQuery("CREATE (f:Forum{Name:$name}) WITH f "+
		"MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id WITH usr, f MERGE (usr)-[:ACTIVE]->(f) "+
		query, dict)
	if err != nil {
		if strings.Contains(err.Error(), "ConstraintValidationFailed") {
			return http.StatusConflict, errors.New("Forum with that name already exists.")
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func SubcribeToForum(token string, forum_name string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr) WHERE ELEMENTID(usr) = %Id "+
		"MATCH (f:Forum) WHERE f.Name = $Name "+
		"MERGE (usr)-[:ACTIVE]->(f)", map[string]any{
		"Id":   user_node.ElementId,
		"Name": forum_name,
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}

func AddPost(token string, forum_name string, post_title string, post_body string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	UUID := GenerateUUID()
	_, err = doQuery("MATCH (f:Forum) WHERE f.Name = $Name "+
		"WITH f MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"AND f IS NOT NULL MERGE (usr)-[:ACTIVE]->(f) WITH usr, f "+
		"CREATE (p:Post {UUID: $UUID, Title: $Title, PostedOn: $TimeStamp})<-[:POSTED_BY]-(usr) "+
		"WITH f, p MERGE (f)-[:HAS_POST]->(p)", map[string]any{
		"Name":      forum_name,
		"Id":        user_node.ElementId,
		"Title":     post_title,
		"UUID":      UUID,
		"TimeStamp": time.Now().Unix(),
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func CommentPost(token string, postUUID string, post_body string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	UUID := GenerateUUID()
	_, err = doQuery("MATCH (post:Post) WHERE post.UUID = $PostUUID "+
		"MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (f:Forum) WHERE (f)-[:HAS_POST]->(post) MERGE (usr)-[:ACTIVE]->(f) WITH usr, post "+
		"CREATE (c:Comment {UUID: $UUID, PostedOn: $TimeStamp})<-[:POSTED_BY]-(usr) "+
		"WITH post, c MERGE (post)-[:HAS_COMMENT]->(c)", map[string]any{
		"PostUUID":  postUUID,
		"Id":        user_node.ElementId,
		"UUID":      UUID,
		"TimeStamp": time.Now().Unix(),
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
