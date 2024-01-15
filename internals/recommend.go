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

type FriendRecommendation struct {
	Username  string
	Interests map[string]bool
	Friends   map[string]bool
	Reason    string
}

func RecommendForums(token string) ([]ForumRecommendation, int, error) {
	rec_slice := make([]ForumRecommendation, 0)

	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}
	rec_by_interests, err := RecommendForumByInterests(user_node.ElementId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rec_by_friends, err := RecommendForumByFriends(user_node.ElementId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	keys := make(map[string]bool, 0)
	for k := range rec_by_interests {
		keys[k] = true
	}
	for k := range rec_by_friends {
		keys[k] = true
	}
	for forum_name := range keys {
		rec_by_int, in_rec_by_int := rec_by_interests[forum_name]
		rec_by_fren, in_rec_by_fren := rec_by_friends[forum_name]
		var rec = ForumRecommendation{}
		if in_rec_by_int && in_rec_by_fren {
			rec = rec_by_int
			rec.Friends = rec_by_fren.Friends
		} else if in_rec_by_int {
			rec = rec_by_int
		} else if in_rec_by_fren {
			rec = rec_by_fren
		}
		rec_slice = append(rec_slice, rec)
	}

	return rec_slice, http.StatusOK, nil
}

func RecommendForumByFriends(user_id string) (map[string]ForumRecommendation, error) {
	result, err := doQuery("MATCH (usr) WHERE ELEMENTID(usr) = $Id MATCH (friends:AccountCredentials) WHERE (usr)-[:FRIEND]-(friends) "+
		"WITH friends, usr MATCH (forums)-[relationship:ACTIVE]-(friends) "+
		"UNWIND relationship as r RETURN STARTNODE(r) AS friend, ENDNODE(r) AS forum ",
		map[string]any{
			"Id": user_id,
		})
	if err != nil {
		return nil, err
	}
	rec_map := make(map[string]ForumRecommendation)
	for _, record := range result.Records {
		forum_node, _ := record.Get("forum")
		friend_node, _ := record.Get("friend")
		forum_name := forum_node.(dbtype.Node).Props["Name"].(string)
		friend_name := friend_node.(dbtype.Node).Props["Username"].(string)
		_, ok := rec_map[forum_name]
		if !ok {
			rec_map[forum_name] = ForumRecommendation{forum_name, make(map[string]bool), make(map[string]bool), ""}
		}
		forum := rec_map[forum_name]
		friends := forum.Friends
		_, ok = friends[friend_name]
		if !ok {
			friends[friend_name] = true
		}
		rec_map[forum_name] = forum
	}

	return rec_map, nil
}

func RecommendForumByInterests(user_id string) (map[string]ForumRecommendation, error) {
	result, err := doQuery("MATCH (usr) WHERE ELEMENTID(usr) = $Id MATCH (interest:Interest) WHERE (usr)-[:INTERESTED_IN]-(interest) "+
		"WITH interest, usr MATCH (forums)-[relationship:HOSTS]->(interest) "+
		"UNWIND relationship as r RETURN STARTNODE(r) AS forum, ENDNODE(r) AS interest ",

		map[string]any{
			"Id": user_id,
		})
	if err != nil {
		return nil, err
	}
	rec_map := make(map[string]ForumRecommendation)
	for _, record := range result.Records {
		forum_node, _ := record.Get("forum")
		interest_node, _ := record.Get("interest")
		forum_name := forum_node.(dbtype.Node).Props["Name"].(string)
		interest_name := interest_node.(dbtype.Node).Props["Name"].(string)
		_, ok := rec_map[forum_name]
		if !ok {
			rec_map[forum_name] = ForumRecommendation{forum_name, make(map[string]bool), make(map[string]bool), ""}
		}
		forum := rec_map[forum_name]
		interests := forum.Interests
		_, ok = interests[interest_name]
		if !ok {
			interests[interest_name] = true
		}
		rec_map[forum_name] = forum
	}
	return rec_map, nil
}

func RecommendFriends(token string) ([]FriendRecommendation, int, error) {
	rec_slice := make([]FriendRecommendation, 0)

	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}
	rec_by_friends, err := RecommendFriendByFriends(user_node.ElementId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rec_by_interests, err := RecommendFriendByInterests(user_node.ElementId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	keys := make([]string, 0)
	for k := range rec_by_interests {
		keys = append(keys, k)
	}
	for k := range rec_by_friends {
		keys = append(keys, k)
	}
	for _, friend_name := range keys {
		rec_by_int, in_rec_by_int := rec_by_interests[friend_name]
		rec_by_fren, in_rec_by_fren := rec_by_friends[friend_name]
		var rec = FriendRecommendation{}
		if in_rec_by_int && in_rec_by_fren {
			rec = rec_by_int
			rec.Friends = rec_by_fren.Friends
		} else if in_rec_by_int {
			rec = rec_by_int
		} else if in_rec_by_fren {
			rec = rec_by_fren
		}
		rec_slice = append(rec_slice, rec)
	}

	return rec_slice, http.StatusOK, nil
}

func RecommendFriendByFriends(user_id string) (map[string]FriendRecommendation, error) {
	result, err := doQuery("MATCH (usr) WHERE ELEMENTID(usr) = $Id "+
		"MATCH friend_path = (usr)-[:FRIEND*2]-(:AccountCredentials) "+
		"WITH nodes(friend_path) AS nodes RETURN nodes[1] AS usr_friend, nodes[2] AS rec_friend",
		map[string]any{
			"Id": user_id,
		})
	if err != nil {
		return nil, err
	}
	rec_map := make(map[string]FriendRecommendation)
	for _, record := range result.Records {
		usr_friend_node, _ := record.Get("usr_friend")
		rec_friend_node, _ := record.Get("rec_friend")
		if err != nil {
			return nil, err
		}
		rec_frend_name := rec_friend_node.(dbtype.Node).Props["Username"].(string)
		usr_friend_name := usr_friend_node.(dbtype.Node).Props["Username"].(string)
		_, ok := rec_map[rec_frend_name]
		if !ok {
			rec_map[rec_frend_name] = FriendRecommendation{rec_frend_name, make(map[string]bool), make(map[string]bool), ""}
		}
		rec_friend := rec_map[rec_frend_name]
		friends := rec_friend.Friends
		_, ok = friends[usr_friend_name]
		if !ok {
			friends[usr_friend_name] = true
		}
		rec_map[rec_frend_name] = rec_friend
	}
	return rec_map, nil
}

func RecommendFriendByInterests(user_id string) (map[string]FriendRecommendation, error) {
	result, err := doQuery("MATCH (usr) WHERE ID(usr) = 0 "+
		"MATCH (interests:Interest) "+
		"WHERE (usr)-[:INTERESTED_IN]-(interests) "+
		"WITH interests, usr "+
		"MATCH (friends)-[relationship:INTERESTED_IN]-(interests) "+
		"WHERE friends <> usr AND NOT (friends)-[:FRIEND]-(usr) "+
		"UNWIND relationship as r RETURN STARTNODE(r) AS friend, ENDNODE(r) AS interest",
		map[string]any{
			"Id": user_id,
		})
	if err != nil {
		return nil, err
	}
	rec_map := make(map[string]FriendRecommendation)
	for _, record := range result.Records {
		interst_node, _ := record.Get("interest")
		friend_node, _ := record.Get("friend")
		interest_name := interst_node.(dbtype.Node).Props["Name"].(string)
		friend_name := friend_node.(dbtype.Node).Props["Username"].(string)
		_, ok := rec_map[friend_name]
		if !ok {
			rec_map[friend_name] = FriendRecommendation{friend_name, make(map[string]bool), make(map[string]bool), ""}
		}
		friend := rec_map[friend_name]
		interests := friend.Interests
		_, ok = interests[interest_name]
		if !ok {
			interests[interest_name] = true
		}
		rec_map[friend_name] = friend
	}
	return rec_map, nil
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
