package internals

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

var expirationTime = 240.0 // 10 dana sad dok radimo, 1 dan u finalnoj verziji

type SessionToken struct {
	Token   string
	Expires int64
}

type AccountCredentials struct {
	Username     string
	PasswordHash string
}

func doQuery(query string, params map[string]any) (*neo4j.EagerResult, error) {
	return neo4j.ExecuteQuery(ctx, driver, query, params, neo4j.EagerResultTransformer)
}

func CreateAccount(username string, password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	if len(username) == 0 {
		return errors.New("Username cannot be empty")
	}

	if len(password) < 6 {
		return errors.New("Password must contain at least 6 characers")
	}

	_, err = doQuery(
		"CREATE (n:AccountCredentials { Username: $Username, PasswordHash: $PasswordHash}) RETURN n",
		map[string]any{
			"Username":     strings.ToLower(username),
			"PasswordHash": hash,
		})
	if err != nil {
		if strings.Contains(err.Error(), "ConstraintValidationFailed") {
			return errors.New("Account with that username already exists.")
		}
		return err
	}
	return nil
}

func GenerateToken() (SessionToken, error) {
	sessionToken := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, sessionToken)
	if err != nil {
		return SessionToken{}, err
	}

	return SessionToken{base64.StdEncoding.EncodeToString(sessionToken), time.Now().Unix() + int64(time.Hour.Seconds()*expirationTime)}, nil
}

func FirstOrDefault(eagerResult *neo4j.EagerResult) (dbtype.Node, bool) {
	if eagerResult == nil {
		return dbtype.Node{}, false
	}
	if eagerResult.Records == nil {
		return dbtype.Node{}, false
	}
	var res dbtype.Node
	res = eagerResult.Records[0].Values[0].(dbtype.Node)

	return res, true

}

func LoginUser(username string, password string) (string, error) {

	sessionToken, err := GenerateToken()
	if err != nil {
		return "", err
	}

	res, err := doQuery("MATCH (ac:AccountCredentials) WHERE ac.Username = $Username RETURN ac", map[string]any{
		"Username": username,
	})
	if err != nil {
		return "", err
	}

	user, ok := FirstOrDefault(res)
	if !ok {
		return "", errors.New("User not found")
	}
	passwordHash := user.Props["PasswordHash"].(string)

	correctPassword, err := argon2id.ComparePasswordAndHash(password, passwordHash)

	if !correctPassword {
		return "", errors.New("Invalid username or password")
	}

	_, err = doQuery(
		"MATCH (ac:AccountCredentials) WHERE ELEMENTID(ac) = $Id "+
			"CREATE (st:SessionToken { Token: $Token, Expires: $Expires })-[:LOGS_IN]->(ac) "+
			"RETURN st",
		map[string]any{
			"Token":   sessionToken.Token,
			"Expires": sessionToken.Expires,
			"Id":      user.ElementId,
		})
	if err != nil {
		return "", err
	}

	return sessionToken.Token, nil
}

func GetUserFromToken(token string) (dbtype.Node, int, error) {
	token_result, err := doQuery("MATCH (st:SessionToken) WHERE st.Token = $Token "+
		"RETURN (st)",
		map[string]any{
			"Token": token,
		})
	if err != nil {
		return dbtype.Node{}, http.StatusInternalServerError, err
	}
	token_node, ok := FirstOrDefault(token_result)
	if !ok {
		return dbtype.Node{}, http.StatusNotFound, errors.New("Token not found")
	}
	exp, _ := token_node.Props["Expires"].(int64)
	if exp < time.Now().Unix() {
		_, err := doQuery("MATCH (st:SessionToken) WHERE ELEMENTID(st) = $Id "+
			"DETACH DELETE st",
			map[string]any{
				"Id": token_node.ElementId,
			})
		if err != nil {
			return dbtype.Node{}, http.StatusNotFound, err
		}

		return dbtype.Node{}, http.StatusUnauthorized, errors.New("Token expired")
	}
	user_result, err := doQuery("MATCH (st:SessionToken) WHERE ELEMENTID(st) = $Id "+
		"MATCH (usr)<-[:LOGS_IN]-(st) "+
		"RETURN usr",
		map[string]any{
			"Id": token_node.ElementId,
		})
	if err != nil {
		return dbtype.Node{}, http.StatusInternalServerError, err
	}
	user_node, ok := FirstOrDefault(user_result)
	if !ok {
		return dbtype.Node{}, http.StatusNotFound, errors.New("User not found")
	}
	return user_node, http.StatusOK, nil
}

func DeleteToken(token string) (int, error) {
	_, err := doQuery("MATCH (st:SessionToken) WHERE st.Token = $Token "+
		"DETACH DELETE st",
		map[string]any{
			"Token": token,
		})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func SendFriendRequest(token string, friendname string) (int, error) {
	friendname = strings.ToLower(friendname)
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (friend:AccountCredentials) WHERE friend.Username = $Username "+
		"AND NOT (usr)-[:FRIEND]-(friend) "+
		"AND NOT (usr)-[:REQUESTS_FRIENDSHIP]-(friend) "+
		"AND usr <> friend "+
		"MERGE (usr)-[:REQUESTS_FRIENDSHIP]->(friend)", map[string]any{
		"Id":       user_node.ElementId,
		"Username": friendname,
	})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func AcceptRequest(token string, friendname string) (int, error) {
	friendname = strings.ToLower(friendname)
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (friend:AccountCredentials) WHERE friend.Username = $Username "+
		"MATCH (friend)-[r:REQUESTS_FRIENDSHIP]->(usr) "+
		"DELETE r "+
		"MERGE (usr)-[:FRIEND]->(friend) ", map[string]any{
		"Id":       user_node.ElementId,
		"Username": friendname,
	})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil

}

func DeclineRequest(token string, friendname string) (int, error) {
	friendname = strings.ToLower(friendname)
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (friend:AccountCredentials) WHERE friend.Username = $Username "+
		"MATCH (friend)-[r:REQUESTS_FRIENDSHIP]-(usr) "+
		"DELETE r ",
		map[string]any{
			"Id":       user_node.ElementId,
			"Username": friendname,
		})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func Unfriend(token string, friendname string) (int, error) {
	friendname = strings.ToLower(friendname)
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (friend:AccountCredentials) WHERE friend.Username = $Username "+
		"MATCH (friend)-[r:FRIEND]-(usr) "+
		"DELETE r ",
		map[string]any{
			"Id":       user_node.ElementId,
			"Username": friendname,
		})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func GetFriends(token string) ([]string, int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}

	result, err := doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (friend)-[:FRIEND]-(usr) "+
		"RETURN friend.Username AS friend", map[string]any{
		"Id": user_node.ElementId,
	})
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	friend_list := make([]string, 0)
	for _, record := range result.Records {
		friend, _ := record.Get("friend")
		friend_list = append(friend_list, friend.(string))
	}

	return friend_list, http.StatusOK, nil
}

func GetFriendRequests(token string) ([]string, int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}

	result, err := doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (friend)-[:REQUESTS_FRIENDSHIP]-(usr) "+
		"RETURN friend.Username AS friend", map[string]any{
		"Id": user_node.ElementId,
	})
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	friend_list := make([]string, 0)
	for _, record := range result.Records {
		friend, _ := record.Get("friend")
		friend_list = append(friend_list, friend.(string))
	}

	return friend_list, http.StatusOK, nil
}
