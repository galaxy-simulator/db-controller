package backend

import (
	"database/sql"
	"fmt"
	"git.darknebu.la/GalaxySimulator/structs"
	"log"
	"strconv"
)

// getBoxWidth gets the width of the box from the node width the given id
func getBoxWidth(nodeID int64) float64 {
	var boxWidth float64

	query := fmt.Sprintf("SELECT box_width FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&boxWidth)
	if err != nil {
		log.Fatalf("[ E ] getBoxWidth query: %v\n\t\t\t query: %s\n", err, query)
	}

	return boxWidth
}

// getTimestepNode gets the timestep of the current node
func getTimestepNode(nodeID int64) int64 {
	var timestep int64

	query := fmt.Sprintf("SELECT timestep FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&timestep)
	if err != nil {
		log.Fatalf("[ E ] getTimeStep query: %v\n\t\t\t query: %s\n", err, query)
	}

	return timestep
}

// getBoxWidth gets the center of the box from the node width the given id
func getBoxCenter(nodeID int64) []float64 {
	var boxCenterX, boxCenterY []uint8

	query := fmt.Sprintf("SELECT box_center[1], box_center[2] FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&boxCenterX, &boxCenterY)
	if err != nil {
		log.Fatalf("[ E ] getBoxCenter query: %v\n\t\t\t query: %s\n", err, query)
	}

	x, parseErr := strconv.ParseFloat(string(boxCenterX), 64)
	y, parseErr := strconv.ParseFloat(string(boxCenterX), 64)

	if parseErr != nil {
		log.Fatalf("[ E ] parse boxCenter: %v\n\t\t\t query: %s\n", err, query)
		log.Fatalf("[ E ] parse boxCenter: (%f, %f)\n", x, y)
	}

	boxCenterFloat := []float64{x, y}

	return boxCenterFloat
}

// getMaxTimestep gets the maximal timestep from the nodes table
func getMaxTimestep() float64 {
	var maxTimestep float64

	query := fmt.Sprintf("SELECT max(timestep) FROM nodes")
	err := db.QueryRow(query).Scan(&maxTimestep)
	if err != nil {
		log.Fatalf("[ E ] getMaxTimestep query: %v\n\t\t\t query: %s\n", err, query)
	}

	return maxTimestep
}

// getStarID returns the id of the star inside of the node with the given ID
func getStarID(nodeID int64) int64 {
	// get the star id from the node
	var starID int64
	query := fmt.Sprintf("SELECT star_id FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&starID)
	if err != nil {
		log.Fatalf("[ E ] getStarID id query: %v\n\t\t\t query: %s\n", err, query)
	}

	return starID
}

// getNodeDepth returns the depth of the given node in the tree
func getNodeDepth(nodeID int64) int64 {
	// build the query
	query := fmt.Sprintf("SELECT depth FROM nodes WHERE node_id=%d", nodeID)

	var depth int64

	// Execute the query
	err := db.QueryRow(query).Scan(&depth)
	if err != nil {
		log.Fatalf("[ E ] getNodeDepth query: %v \n\t\t\t query: %s\n", err, query)
	}

	return depth
}

// quadrant returns the quadrant into which the given star belongs
func quadrant(star structs.Star2D, nodeID int64) int64 {
	// get the center of the node the star is in
	center := getBoxCenter(nodeID)
	centerX := center[0]
	centerY := center[1]

	if star.C.X > centerX {
		if star.C.Y > centerY {
			// North East condition
			return 1
		}
		// South East condition
		return 3
	}

	if star.C.Y > centerY {
		// North West condition
		return 0
	}
	// South West condition
	return 2
}

// getQuadrantNodeID returns the id of the requested child-node
// Example: if a parent has four children and quadrant 0 is requested, the function returns the north east child id
func getQuadrantNodeID(parentNodeID int64, quadrant int64) int64 {
	var a, b, c, d []uint8

	// get the star from the stars table
	query := fmt.Sprintf("SELECT subnode[1], subnode[2], subnode[3], subnode[4] FROM nodes WHERE node_id=%d", parentNodeID)
	err := db.QueryRow(query).Scan(&a, &b, &c, &d)
	if err != nil {
		log.Fatalf("[ E ] getQuadrantNodeID star query: %v \n\t\t\tquery: %s\n", err, query)
	}

	returnA, _ := strconv.ParseInt(string(a), 10, 64)
	returnB, _ := strconv.ParseInt(string(b), 10, 64)
	returnC, _ := strconv.ParseInt(string(c), 10, 64)
	returnD, _ := strconv.ParseInt(string(d), 10, 64)

	switch quadrant {
	case 0:
		return returnA
	case 1:
		return returnB
	case 2:
		return returnC
	case 3:
		return returnD
	}

	return -1
}

// GetStar returns the star with the given ID from the stars table
func GetStar(starID int64) structs.Star2D {
	var x, y, vx, vy, m float64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT x, y, vx, vy, m FROM stars WHERE star_id=%d", starID)
	err := db.QueryRow(query).Scan(&x, &y, &vx, &vy, &m)
	if err != nil {
		log.Fatalf("[ E ] GetStar query: %v \n\t\t\tquery: %s\n", err, query)
	}

	star := structs.Star2D{
		C: structs.Vec2{
			X: x,
			Y: y,
		},
		V: structs.Vec2{
			X: vx,
			Y: vy,
		},
		M: m,
	}

	return star
}

// GetStarIDTimestep returns the timestep the given starID is currently inside of
func GetStarIDTimestep(starID int64) int64 {
	var timestep int64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT timestep FROM nodes WHERE star_id=%d", starID)
	err := db.QueryRow(query).Scan(&timestep)
	if err != nil {
		log.Fatalf("[ E ] GetStar query: %v \n\t\t\tquery: %s\n", err, query)
	}

	return timestep
}

// getStarMass returns the mass if the star with the given ID
func getStarMass(starID int64) float64 {
	var mass float64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT m FROM stars WHERE star_id=%d", starID)
	err := db.QueryRow(query).Scan(&mass)
	if err != nil {
		log.Fatalf("[ E ] getStarMass query: %v \n\t\t\tquery: %s\n", err, query)
	}

	return mass
}

// getNodeTotalMass returns the total mass of the node with the given ID and its children
func getNodeTotalMass(nodeID int64) float64 {
	var mass float64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT total_mass FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&mass)
	if err != nil {
		log.Fatalf("[ E ] getStarMass query: %v \n\t\t\tquery: %s\n", err, query)
	}

	return mass
}

// GetListOfStarsGo returns the list of stars in go struct format
func GetListOfStarsGo(database *sql.DB) []structs.Star2D {
	db = database
	// build the query
	query := fmt.Sprintf("SELECT * FROM stars")

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] removeStarFromNode query: %v\n\t\t\t query: %s\n", err, query)
	}

	var starList []structs.Star2D

	// iterate over the returned rows
	for rows.Next() {

		var starID int64
		var x, y, vx, vy, m float64
		scanErr := rows.Scan(&starID, &x, &y, &vx, &vy, &m)
		if scanErr != nil {
			log.Fatalf("[ E ] scan error: %v", scanErr)
		}

		star := structs.Star2D{
			C: structs.Vec2{
				X: x,
				Y: y,
			},
			V: structs.Vec2{
				X: vx,
				Y: vy,
			},
			M: m,
		}

		starList = append(starList, star)
	}

	return starList
}

// GetListOfStarIDs returns a list of all star ids in the stars table
func GetListOfStarIDs(db *sql.DB) []int64 {
	// build the query
	query := fmt.Sprintf("SELECT star_id FROM stars")

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] GetListOfStarIDs query: %v\n\t\t\t query: %s\n", err, query)
	}

	var starIDList []int64

	// iterate over the returned rows
	for rows.Next() {

		var starID int64
		scanErr := rows.Scan(&starID)
		if scanErr != nil {
			log.Fatalf("[ E ] scan error: %v", scanErr)
		}

		starIDList = append(starIDList, starID)
	}

	return starIDList
}

// GetListOfStarIDs returns a list of all star ids in the stars table with the given timestep
func GetListOfStarIDsTimestep(db *sql.DB, timestep int64) []int64 {
	// build the query
	query := fmt.Sprintf("SELECT star_id FROM nodes WHERE star_id<>0 AND timestep=%d", timestep)

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] GetListOfStarIDsTimestep query: %v\n\t\t\t query: %s\n", err, query)
	}

	var starIDList []int64

	// iterate over the returned rows
	for rows.Next() {

		var starID int64
		scanErr := rows.Scan(&starID)
		if scanErr != nil {
			log.Fatalf("[ E ] scan error: %v", scanErr)
		}

		starIDList = append(starIDList, starID)
	}

	return starIDList
}

// GetListOfStarsCsv returns an array of strings containing the coordinates of all the stars in the stars table
func GetListOfStarsCsv(db *sql.DB) []string {
	// build the query
	query := fmt.Sprintf("SELECT * FROM stars")

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] getListOfStarsCsv query: %v\n\t\t\t query: %s\n", err, query)
	}

	var starList []string

	// iterate over the returned rows
	for rows.Next() {

		var starID int64
		var x, y, vx, vy, m float64
		scanErr := rows.Scan(&starID, &x, &y, &vx, &vy, &m)
		if scanErr != nil {
			log.Fatalf("[ E ] scan error: %v", scanErr)
		}

		row := fmt.Sprintf("%d, %f, %f, %f, %f, %f", starID, x, y, vx, vy, m)
		starList = append(starList, row)
	}

	return starList
}

// GetListOfStarsTreeCsv returns an array of strings containing the coordinates of all the stars in the given tree
func GetListOfStarsTree(database *sql.DB, treeindex int64) []structs.Star2D {
	db = database

	// build the query
	query := fmt.Sprintf("SELECT * FROM stars WHERE star_id IN(SELECT star_id FROM nodes WHERE timestep=%d)", treeindex)

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] removeStarFromNode query: %v\n\t\t\t query: %s\n", err, query)
	}

	var starList []structs.Star2D

	// iterate over the returned rows
	for rows.Next() {

		var starID int64
		var x, y, vx, vy, m float64
		scanErr := rows.Scan(&starID, &x, &y, &vx, &vy, &m)
		if scanErr != nil {
			log.Fatalf("[ E ] scan error: %v", scanErr)
		}

		star := structs.Star2D{
			C: structs.Vec2{
				X: x,
				Y: y,
			},
			V: structs.Vec2{
				X: vx,
				Y: vy,
			},
			M: m,
		}

		starList = append(starList, star)
	}

	return starList
}

// getRootNodeID gets a tree index and returns the nodeID of its root node
func getRootNodeID(index int64) int64 {
	var nodeID int64

	log.Printf("Preparing query with the root id %d", index)
	query := fmt.Sprintf("SELECT node_id FROM nodes WHERE root_id=%d", index)
	log.Printf("Sending query")
	err := db.QueryRow(query).Scan(&nodeID)
	if err != nil {
		log.Fatalf("[ E ] getRootNodeID query: %v\n\t\t\t query: %s\n", err, query)
	}
	log.Printf("Done Sending query")

	return nodeID
}

// getCenterOfMass returns the center of mass of the given nodeID
func getCenterOfMass(nodeID int64) structs.Vec2 {

	var CenterOfMass [2]float64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT center_of_mass[1], center_of_mass[2] FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&CenterOfMass[0], &CenterOfMass[1])
	if err != nil {
		log.Fatalf("[ E ] getCenterOfMass query: %v \n\t\t\tquery: %s\n", err, query)
	}

	return structs.Vec2{X: CenterOfMass[0], Y: CenterOfMass[1]}
}

// getStarCoordinates gets the star coordinates of a star using a given nodeID.
// It returns a vector describing the coordinates
func getStarCoordinates(nodeID int64) structs.Vec2 {
	var Coordinates [2]float64

	starID := getStarID(nodeID)

	// get the star from the stars table
	query := fmt.Sprintf("SELECT x, y FROM stars WHERE star_id=%d", starID)
	err := db.QueryRow(query).Scan(&Coordinates[0], &Coordinates[1])
	if err != nil {
		log.Fatalf("[ E ] getStarCoordinates query: %v \n\t\t\tquery: %s\n", err, query)
	}

	fmt.Printf("%v\n", Coordinates)

	return structs.Vec2{X: Coordinates[0], Y: Coordinates[1]}
}

// getNodeCenterOfMass returns the center of mass of the node with the given ID
func getNodeCenterOfMass(nodeID int64) structs.Vec2 {
	var Coordinates [2]float64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT center_of_mass[1], center_of_mass[2] FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&Coordinates[0], &Coordinates[1])
	if err != nil {
		log.Fatalf("[ E ] getNodeCenterOfMass query: %v \n\t\t\tquery: %s\n", err, query)
	}

	return structs.Vec2{X: Coordinates[0], Y: Coordinates[1]}
}

// getSubtreeIDs returns the id of the subtrees of the nodeID
func getSubtreeIDs(nodeID int64) [4]int64 {

	var subtreeIDs [4]int64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT subnode[1], subnode[2], subnode[3], subnode[4] FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&subtreeIDs[0], &subtreeIDs[1], &subtreeIDs[2], &subtreeIDs[3])
	if err != nil {
		log.Fatalf("[ E ] getSubtreeIDs query: %v \n\t\t\tquery: %s\n", err, query)
	}

	return subtreeIDs
}
