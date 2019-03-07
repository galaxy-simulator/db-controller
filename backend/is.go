package backend

import (
	"fmt"
	"log"
)

// isLeaf returns true if the node with the given id is a leaf
func isLeaf(nodeID int64) bool {
	var isLeaf bool

	query := fmt.Sprintf("SELECT COALESCE(isleaf, FALSE) FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&isLeaf)
	if err != nil {
		log.Fatalf("[ E ] isLeaf query: %v\n\t\t\t query: %s\n", err, query)
	}

	if isLeaf == true {
		return true
	}

	return false
}
