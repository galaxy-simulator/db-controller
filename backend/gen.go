package backend

import (
	"database/sql"
	"fmt"
	"log"
)

// genForestTree generates a forest representation of the tree with the given index
func GenForestTree(database *sql.DB, index int64) string {
	db = database
	rootNodeID := getRootNodeID(index)
	return genForestTreeNode(rootNodeID)
}

// genForestTreeNodes returns a sub-representation of a given node in forest format
func genForestTreeNode(nodeID int64) string {
	var returnString string

	// get the subnode ids
	var subnode [4]int64

	query := fmt.Sprintf("SELECT subnode[1], subnode[2], subnode[3], subnode[4] FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&subnode[0], &subnode[1], &subnode[2], &subnode[3])
	if err != nil {
		log.Fatalf("[ E ] updateTotalMassNode query: %v\n\t\t\t query: %s\n", err, query)
	}

	returnString += "["

	// iterate over all subnodes updating their total masses
	for _, subnodeID := range subnode {
		if subnodeID != 0 {
			centerOfMass := getCenterOfMass(nodeID)
			mass := getNodeTotalMass(nodeID)
			returnString += fmt.Sprintf("%.0f %.0f %.0f", centerOfMass.X, centerOfMass.Y, mass)
			returnString += genForestTreeNode(subnodeID)
		} else {
			if getStarID(nodeID) != 0 {
				coords := getStarCoordinates(nodeID)
				starID := getStarID(nodeID)
				mass := getStarMass(starID)
				returnString += fmt.Sprintf("[%.0f %.0f %.0f]", coords.X, coords.Y, mass)
			} else {
				returnString += fmt.Sprintf("[0 0]")
			}
			// break, this stops a star from being counted multiple (4) times
			break
		}
	}

	returnString += "]"

	return returnString
}
