package backend

import (
	"fmt"
	"log"
)

// containsStar returns true if the node with the given id contains a star and returns false if not.
func containsStar(id int64) bool {
	var starID int64

	query := fmt.Sprintf("SELECT star_id FROM nodes WHERE node_id=%d", id)
	err := db.QueryRow(query).Scan(&starID)
	if err != nil {
		log.Fatalf("[ E ] containsStar query: %v\n\t\t\t query: %s\n", err, query)
	}

	if starID != 0 {
		return true
	}

	return false
}
