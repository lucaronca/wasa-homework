package repositories

import "strings"

type Relation func(entity string) string

func queryBuilder(entity string, relations ...Relation) string {
	relationsStrings := make([]string, len(relations))
	alreadyAddedWhereRelation := false
	for i, relation := range relations {
		relString := relation(entity)
		if !strings.Contains(relString, "JOIN") && strings.Contains(relString, "WHERE") {
			if alreadyAddedWhereRelation {
				relString = strings.Replace(relString, "WHERE", "AND", 1)
			} else {
				alreadyAddedWhereRelation = true
			}
		}
		relationsStrings[i] = relString
	}

	return strings.Join(relationsStrings, " ")
}

var dateLayout = "2006-01-02T15:04:05-0700"
