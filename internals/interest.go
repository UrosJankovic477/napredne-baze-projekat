package internals

import "net/http"

type Interest struct {
	Name string
}

func AddInterest(token string, categry string, interest string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MERGE (interest:Interest{Name:LOWER($name)}) "+
		"MERGE (usr)-[:INTERESTED_IN]->(interest) "+
		"MERGE (cat:Interest{Name:LOWER($category)}) "+
		"MERGE (cat)-[:IS_CATEGORY_OF]->(interest)",
		map[string]any{
			"Id":       user_node.ElementId,
			"name":     interest,
			"category": categry,
		})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
