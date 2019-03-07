package backend

import (
	"database/sql"
	"fmt"
	"git.darknebu.la/GalaxySimulator/structs"
	"log"
)

// UpdateTotalMass gets a tree index and returns the nodeID of the trees root node
func UpdateTotalMass(database *sql.DB, index int64) {
	db = database
	rootNodeID := getRootNodeID(index)
	log.Printf("RootID: %d", rootNodeID)
	updateTotalMassNode(rootNodeID)
}

// updateTotalMassNode updates the total mass of the given node
func updateTotalMassNode(nodeID int64) float64 {
	var totalmass float64

	// get the subnode ids
	var subnode [4]int64

	query := fmt.Sprintf("SELECT subnode[1], subnode[2], subnode[3], subnode[4] FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&subnode[0], &subnode[1], &subnode[2], &subnode[3])
	if err != nil {
		log.Fatalf("[ E ] updateTotalMassNode query: %v\n\t\t\t query: %s\n", err, query)
	}
	// TODO: implement the getSubtreeIDs(nodeID) []int64 {...} function
	// iterate over all subnodes updating their total masses
	for _, subnodeID := range subnode {
		fmt.Println("----------------------------")
		fmt.Printf("SubdnodeID: %d\n", subnodeID)
		if subnodeID != 0 {
			totalmass += updateTotalMassNode(subnodeID)
		} else {
			// get the starID for getting the star mass
			starID := getStarID(nodeID)
			fmt.Printf("StarID: %d\n", starID)
			if starID != 0 {
				mass := getStarMass(starID)
				log.Printf("starID=%d \t mass: %f", starID, mass)
				totalmass += mass
			}

			// break, this stops a star from being counted multiple (4) times
			break
		}
		fmt.Println("----------------------------")
	}

	query = fmt.Sprintf("UPDATE nodes SET total_mass=%f WHERE node_id=%d", totalmass, nodeID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] insert total_mass query: %v\n\t\t\t query: %s\n", err, query)
	}

	fmt.Printf("nodeID: %d \t totalMass: %f\n", nodeID, totalmass)

	return totalmass
}

// UpdateCenterOfMass recursively updates the center of mass of all the nodes starting at the node with the given
// root index
func UpdateCenterOfMass(database *sql.DB, index int64) {
	db = database
	rootNodeID := getRootNodeID(index)
	log.Printf("RootID: %d", rootNodeID)
	updateCenterOfMassNode(rootNodeID)
}

// updateCenterOfMassNode updates the center of mass of the node with the given nodeID recursively
// center of mass := ((x_1 * m) + (x_2 * m) + ... + (x_n * m)) / m
func updateCenterOfMassNode(nodeID int64) structs.Vec2 {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	var centerOfMass structs.Vec2

	// get the subnode ids
	var subnode [4]int64
	var starID int64

	query := fmt.Sprintf("SELECT subnode[1], subnode[2], subnode[3], subnode[4], star_id FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&subnode[0], &subnode[1], &subnode[2], &subnode[3], &starID)
	if err != nil {
		log.Fatalf("[ E ] updateCenterOfMassNode query: %v\n\t\t\t query: %s\n", err, query)
	}

	// if the nodes does not contain a star but has children, update the center of mass
	if subnode != ([4]int64{0, 0, 0, 0}) {
		log.Println("[   ] recursing deeper")

		// define variables storing the values of the subnodes
		var totalMass float64
		var centerOfMassX float64
		var centerOfMassY float64

		// iterate over all the subnodes and calculate the center of mass of each node
		for _, subnodeID := range subnode {
			subnodeCenterOfMass := updateCenterOfMassNode(subnodeID)

			if subnodeCenterOfMass.X != 0 && subnodeCenterOfMass.Y != 0 {
				fmt.Printf("SubnodeCenterOfMass: (%f, %f)\n", subnodeCenterOfMass.X, subnodeCenterOfMass.Y)
				subnodeMass := getNodeTotalMass(subnodeID)
				totalMass += subnodeMass

				centerOfMassX += subnodeCenterOfMass.X * subnodeMass
				centerOfMassY += subnodeCenterOfMass.Y * subnodeMass
			}
		}

		// calculate the overall center of mass of the subtree
		centerOfMass = structs.Vec2{
			X: centerOfMassX / totalMass,
			Y: centerOfMassY / totalMass,
		}

		// else, use the star as the center of mass (this can be done, because of the rule defining that there
		// can only be one star in a cell)
	} else {
		log.Println("[   ] using the star in the node as the center of mass")
		log.Printf("[   ] NodeID: %v", nodeID)
		starID := getStarID(nodeID)

		if starID == 0 {
			log.Println("[   ] StarID == 0...")
			centerOfMass = structs.Vec2{
				X: 0,
				Y: 0,
			}
		} else {
			log.Printf("[   ] NodeID: %v", starID)
			star := GetStar(starID)
			centerOfMassX := star.C.X
			centerOfMassY := star.C.Y
			centerOfMass = structs.Vec2{
				X: centerOfMassX,
				Y: centerOfMassY,
			}
		}
	}

	// build the query
	query = fmt.Sprintf("UPDATE nodes SET center_of_mass='{%f, %f}' WHERE node_id=%d", centerOfMass.X, centerOfMass.Y, nodeID)

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] update center of mass query: %v\n\t\t\t query: %s\n", err, query)
	}

	fmt.Printf("[   ] CenterOfMass: (%f, %f)\n", centerOfMass.X, centerOfMass.Y)

	return centerOfMass
}

// updateStarForce updates the force acting on the star
func updateStarForce(db *sql.DB, starID int64, force structs.Vec2) structs.Star2D {

	star := GetStar(starID)
	newStar := structs.Star2D{
		structs.Vec2{star.C.X, star.C.Y},
		structs.Vec2{force.X, force.Y},
		star.M,
	}

	// updated the stars Force
	query := fmt.Sprintf("UPDATE stars SET vx=%f, vy=%f WHERE star_id=%d", force.X, force.Y, starID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] updateStarForce query: %v\n\t\t\t query: %s\n", err, query)
	}

	return newStar
}
