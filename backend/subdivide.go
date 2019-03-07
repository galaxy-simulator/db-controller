package backend

import (
	"fmt"
	"log"
)

// subdivide subdivides the given node creating four child nodes
func subdivide(nodeID int64) {
	boxWidth := getBoxWidth(nodeID)
	boxCenter := getBoxCenter(nodeID)
	originalDepth := getNodeDepth(nodeID)
	timestep := getTimestepNode(nodeID)
	log.Printf("Subdividing %d, setting the timestep to %d", nodeID, timestep)

	// calculate the new positions
	newPosX := boxCenter[0] + (boxWidth / 2)
	newPosY := boxCenter[1] + (boxWidth / 2)
	newNegX := boxCenter[0] - (boxWidth / 2)
	newNegY := boxCenter[1] - (boxWidth / 2)
	newWidth := boxWidth / 2

	// create new news with those positions
	newNodeIDA := newNode(newPosX, newPosY, newWidth, originalDepth+1, timestep)
	newNodeIDB := newNode(newPosX, newNegY, newWidth, originalDepth+1, timestep)
	newNodeIDC := newNode(newNegX, newPosY, newWidth, originalDepth+1, timestep)
	newNodeIDD := newNode(newNegX, newNegY, newWidth, originalDepth+1, timestep)

	// Update the subtrees of the parent node

	// build the query
	query := fmt.Sprintf("UPDATE nodes SET subnode='{%d, %d, %d, %d}', isleaf=FALSE, timestep=%d WHERE node_id=%d", newNodeIDA, newNodeIDB, newNodeIDC, newNodeIDD, timestep, nodeID)

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] subdivide query: %v\n\t\t\t query: %s\n", err, query)
	}
}
