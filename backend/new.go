package backend

import (
	"database/sql"
	"fmt"
	"log"
)

// newTree creates a new tree with the given width
func NewTree(database *sql.DB, width float64) {
	db = database

	log.Printf("Creating a new tree with a width of %f", width)

	// get the current max root id
	query := fmt.Sprintf("SELECT COALESCE(max(root_id), 0) FROM nodes")
	var currentMaxRootID int64
	err := db.QueryRow(query).Scan(&currentMaxRootID)
	if err != nil {
		log.Fatalf("[ E ] max root id query: %v\n\t\t\t query: %s\n", err, query)
	}

	// build the query creating a new node
	query = fmt.Sprintf("INSERT INTO nodes (box_width, root_id, box_center, depth, isleaf, timestep) VALUES (%f, %d, '{0, 0}', 0, TRUE, %d)", width, currentMaxRootID+1, currentMaxRootID+1)

	// execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] insert new node query: %v\n\t\t\t query: %s\n", err, query)
	}
}

// newNode Inserts a new node into the database with the given parameters
func newNode(x float64, y float64, width float64, depth int64, timestep int64) int64 {
	// build the query creating a new node
	query := fmt.Sprintf("INSERT INTO nodes (box_center, box_width, depth, isleaf, timestep) VALUES ('{%f, %f}', %f, %d, TRUE, %d) RETURNING node_id", x, y, width, depth, timestep)

	var nodeID int64

	// execute the query
	err := db.QueryRow(query).Scan(&nodeID)
	if err != nil {
		log.Fatalf("[ E ] newNode query: %v\n\t\t\t query: %s\n", err, query)
	}

	return nodeID
}
