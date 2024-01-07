package internals

import (
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type ForumRecommendation struct {
	Name      string
	Interests map[string]bool
	Friends   map[string]bool
	Reason    string
}

func RecommendForums(token string) ([]ForumRecommendation, int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}
	rec_by_interests, err := RecommendForumByFriends(user_node.ElementId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	//for _, v := range rec_by_interests {
	//
	//}
	return rec_by_interests, http.StatusNotImplemented, nil
}

func RecommendForumByFriends(user_id string) ([]ForumRecommendation, error) {
	result, err := doQuery("MATCH (usr) WHERE ELEMENTID(usr) = $Id MATCH (friends:AccountCredentials) WHERE (usr)-[:FRIEND]-(friends) "+
		"WITH friends, usr MATCH (friends)-[relationship:ACTIVE]-(forums) "+
		"UNWIND relationship as r WITH STARTNODE(r) AS friend, ENDNODE(r) AS forum "+
		"MATCH (interests:Interest) WHERE (forum)-[:HOSTS]->(interests) RETURN friend, forum, interests",
		map[string]any{
			"Id": user_id,
		})
	if err != nil {
		return nil, err
	}
	recomendations := make([]ForumRecommendation, 0)
	recommendation := make(map[string]ForumRecommendation)
	for _, record := range result.Records {
		friend_node := record.Values[0].(dbtype.Node)
		forum_node := record.Values[1].(dbtype.Node)
		interest_node := record.Values[2].(dbtype.Node)
		friend_name := friend_node.Props["Username"].(string)
		forum_name := forum_node.Props["Name"].(string)
		interest_name := interest_node.Props["Name"].(string)
		_, ok := recommendation[forum_name]
		if !ok {
			recommendation[forum_name] = ForumRecommendation{forum_name, make(map[string]bool), make(map[string]bool), ""}
		}
		forum := recommendation[forum_name]
		interests := forum.Interests
		_, ok = interests[interest_name]
		if !ok {
			interests[interest_name] = true
		}
		friends := forum.Friends
		_, ok = friends[friend_name]
		if !ok {
			friends[friend_name] = true
		}
		recommendation[forum_name] = forum
	}
	for _, v := range recommendation {
		recomendations = append(recomendations, v)
	}

	return recomendations, nil
}

func RecommendForumByInterests(user_id string) ([]ForumRecommendation, error) {
	result, err := doQuery("MATCH (usr) WHERE ELEMENTID(usr) = $Id MATCH (interest:Interest) WHERE (usr)-[:INTERESTED_IN]-(interest) "+
		"WITH interest, usr MATCH (forums)-[relationship:HOSTS]->(interest) "+
		"UNWIND relationship as r RETURN STARTNODE(r) AS forum, ENDNODE(r) AS interest ",

		map[string]any{
			"Id": user_id,
		})
	if err != nil {
		return nil, err
	}
	recomendations := make([]ForumRecommendation, 0)
	recommendation := make(map[string]ForumRecommendation)
	for _, record := range result.Records {
		forum_node := record.Values[0].(dbtype.Node)
		interest_node := record.Values[1].(dbtype.Node)
		forum_name := forum_node.Props["Name"].(string)
		interest_name := interest_node.Props["Name"].(string)
		_, ok := recommendation[forum_name]
		if !ok {
			recommendation[forum_name] = ForumRecommendation{forum_name, make(map[string]bool), make(map[string]bool), ""}
		}
		forum := recommendation[forum_name]
		interests := forum.Interests
		_, ok = interests[interest_name]
		if !ok {
			interests[interest_name] = true
		}
		recommendation[forum_name] = forum
	}
	for _, v := range recommendation {
		recomendations = append(recomendations, v)
	}
	return recomendations, nil
}

func ToList(eagerResult *neo4j.EagerResult) ([]dbtype.Node, bool) {
	if eagerResult == nil {
		return nil, false
	}
	if eagerResult.Records == nil {
		return nil, false
	}
	var res []dbtype.Node
	res = make([]dbtype.Node, 0)
	for _, record := range eagerResult.Records {
		for _, value := range record.Values {
			res = append(res, value.(dbtype.Node))
		}
	}
	return res, true
}
