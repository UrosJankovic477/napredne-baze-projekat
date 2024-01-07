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
	query := "MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id " +
		"MERGE (interest:Interest{Name:toLower($name)}) " +
		"MERGE (usr)-[:INTERESTED_IN]->(interest) "
	if categry != "" {
		query += "MERGE (cat:Interest{Name:toLower($category)}) " +
			"MERGE (cat)-[:IS_CATEGORY_OF]->(interest)"
	}

	_, err = doQuery(query,
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

func RemoveInterest(token string, interest string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}
	_, err = doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH r WHERE (usr)-[r:INTERESTED_IN]->(interest) "+
		"DELETE r",
		map[string]any{
			"Id":   user_node.ElementId,
			"name": interest,
		})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
