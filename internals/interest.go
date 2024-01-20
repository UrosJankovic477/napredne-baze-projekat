package internals

import (
	"errors"
	"net/http"
)

type Interest struct {
	Name string
}

func AddInterest(token string, categry string, interest string) (int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return status, err
	}

	if len(interest) == 0 {
		return status, errors.New("Interest cannot be empty")
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
		"MATCH (interest) WHERE interest.Name = $interest "+
		"MATCH (usr)-[r:INTERESTED_IN]->(interest) "+
		"DELETE r",
		map[string]any{
			"Id":       user_node.ElementId,
			"interest": interest,
		})
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func GetInterests(token string) ([]string, int, error) {
	user_node, status, err := GetUserFromToken(token)
	if err != nil {
		return nil, status, err
	}
	result, err := doQuery("MATCH (usr:AccountCredentials) WHERE ELEMENTID(usr) = $Id "+
		"MATCH (interest) WHERE (usr)-[:INTERESTED_IN]->(interest) "+
		"RETURN interest.Name AS interest",
		map[string]any{
			"Id": user_node.ElementId,
		})
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	interest_list := make([]string, 0)
	for _, record := range result.Records {
		interest, _ := record.Get("interest")
		interest_list = append(interest_list, interest.(string))
	}

	return interest_list, http.StatusOK, nil
}
