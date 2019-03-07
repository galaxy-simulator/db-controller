package backend

import (
	"database/sql"
	"fmt"
	"log"
)

// deleteAll Stars deletes all the rows in the stars table
func DeleteAllStars(database *sql.DB) {
	db = database
	// build the query creating a new node
	query := "DELETE FROM stars WHERE TRUE"

	// execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] deleteAllStars query: %v\n\t\t\t query: %s\n", err, query)
	}
}

// deleteAll Stars deletes all the rows in the nodes table
func DeleteAllNodes(database *sql.DB) {
	db = database
	// build the query creating a new node
	query := "DELETE FROM nodes WHERE TRUE"

	// execute the query
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("[ E ] deleteAllStars query: %v\n\t\t\t query: %s\n", err, query)
	}
}

// removeStarFromNode removes the star from the node with the given ID
func removeStarFromNode(nodeID int64) {
	// build the query
	query := fmt.Sprintf("UPDATE nodes SET star_id=0 WHERE node_id=%d", nodeID)

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] removeStarFromNode query: %v\n\t\t\t query: %s\n", err, query)
	}
}
