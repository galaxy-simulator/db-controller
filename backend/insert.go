package backend

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"git.darknebu.la/GalaxySimulator/structs"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// InsertStar inserts the given star into the stars table and the nodes table tree
func InsertStar(database *sql.DB, star structs.Star2D, index int64) int64 {
	db = database
	start := time.Now()

	log.Printf("Inserting the star %v into the tree with the index %d", star, index)

	// insert the star into the stars table
	starID := insertIntoStars(star)

	// get the root node id
	query := fmt.Sprintf("select case when exists (select node_id from nodes where root_id=%d) then (select node_id from nodes where root_id=%d) else -1 end;", index, index)
	var id int64
	err := db.QueryRow(query).Scan(&id)

	// if there are no rows in the result set, create a new tree
	if err != nil {
		log.Fatalf("[ E ] Get root node id query: %v\n\t\t\t query: %s\n", err, query)
	}

	if id == -1 {
		NewTree(db, 1000)
		id = getRootNodeID(index)
	}

	log.Printf("Node id of the root node %d: %d", id, index)

	// insert the star into the tree (using it's ID) starting at the root
	insertIntoTree(starID, id)
	elapsedTime := time.Since(start)
	log.Printf("\t\t\t\t\t %s", elapsedTime)
	return starID
}

// insertIntoStars inserts the given star into the stars table
func insertIntoStars(star structs.Star2D) int64 {
	// unpack the star
	x := star.C.X
	y := star.C.Y
	vx := star.V.X
	vy := star.V.Y
	m := star.M

	// build the request query
	query := fmt.Sprintf("INSERT INTO stars (x, y, vx, vy, m) VALUES (%f, %f, %f, %f, %f) RETURNING star_id", x, y, vx, vy, m)

	// execute the query
	var starID int64
	err := db.QueryRow(query).Scan(&starID)
	if err != nil {
		log.Fatalf("[ E ] insert query: %v\n\t\t\t query: %s\n", err, query)
	}

	return starID
}

// insert into tree inserts the given star into the tree starting at the node with the given node id
func insertIntoTree(starID int64, nodeID int64) {
	//starRaw := GetStar(starID)
	//nodeCenter := getBoxCenter(nodeID)
	//nodeWidth := getBoxWidth(nodeID)
	//log.Printf("[   ] \t Inserting star %v into the node (c: %v, w: %v)", starRaw, nodeCenter, nodeWidth)

	// There exist four cases:
	//                    | Contains a Star | Does not Contain a Star |
	// ------------------ + --------------- + ----------------------- +
	// Node is a Leaf     | Impossible      | insert into node        |
	//                    |                 | subdivide               |
	// ------------------ + --------------- + ----------------------- +
	// Node is not a Leaf | insert preexist | insert into the subtree |
	//                    | insert new      |                         |
	// ------------------ + --------------- + ----------------------- +

	// get the node with the given nodeID
	// find out if the node contains a star or not
	containsStar := containsStar(nodeID)

	// find out if the node is a leaf
	isLeaf := isLeaf(nodeID)

	// if the node is a leaf and contains a star
	// subdivide the tree
	// insert the preexisting star into the correct subtree
	// insert the new star into the subtree
	if isLeaf == true && containsStar == true {
		//log.Printf("Case 1, \t %v \t %v", nodeWidth, nodeCenter)
		subdivide(nodeID)
		//tree := printTree(nodeID)

		// Stage 1: Inserting the blocking star
		blockingStarID := getStarID(nodeID)                               // get the id of the star blocking the node
		blockingStar := GetStar(blockingStarID)                           // get the actual star
		blockingStarQuadrant := quadrant(blockingStar, nodeID)            // find out in which quadrant it belongs
		quadrantNodeID := getQuadrantNodeID(nodeID, blockingStarQuadrant) // get the nodeID of that quadrant
		insertIntoTree(blockingStarID, quadrantNodeID)                    // insert the star into that node
		removeStarFromNode(nodeID)                                        // remove the blocking star from the node it was blocking

		// Stage 1: Inserting the actual star
		star := GetStar(starID)                                  // get the actual star
		starQuadrant := quadrant(star, nodeID)                   // find out in which quadrant it belongs
		quadrantNodeID = getQuadrantNodeID(nodeID, starQuadrant) // get the nodeID of that quadrant
		insertIntoTree(starID, nodeID)
	}

	// if the node is a leaf and does not contain a star
	// insert the star into the node and subdivide it
	if isLeaf == true && containsStar == false {
		//log.Printf("Case 2, \t %v \t %v", nodeWidth, nodeCenter)
		directInsert(starID, nodeID)
	}

	// if the node is not a leaf and contains a star
	// insert the preexisting star into the correct subtree
	// insert the new star into the subtree
	if isLeaf == false && containsStar == true {
		//log.Printf("Case 3, \t %v \t %v", nodeWidth, nodeCenter)
		// Stage 1: Inserting the blocking star
		blockingStarID := getStarID(nodeID)                               // get the id of the star blocking the node
		blockingStar := GetStar(blockingStarID)                           // get the actual star
		blockingStarQuadrant := quadrant(blockingStar, nodeID)            // find out in which quadrant it belongs
		quadrantNodeID := getQuadrantNodeID(nodeID, blockingStarQuadrant) // get the nodeID of that quadrant
		insertIntoTree(blockingStarID, quadrantNodeID)                    // insert the star into that node
		removeStarFromNode(nodeID)                                        // remove the blocking star from the node it was blocking

		// Stage 1: Inserting the actual star
		star := GetStar(blockingStarID)                          // get the actual star
		starQuadrant := quadrant(star, nodeID)                   // find out in which quadrant it belongs
		quadrantNodeID = getQuadrantNodeID(nodeID, starQuadrant) // get the nodeID of that quadrant
		insertIntoTree(starID, nodeID)
	}

	// if the node is not a leaf and does not contain a star
	// insert the new star into the according subtree
	if isLeaf == false && containsStar == false {
		//log.Printf("Case 4, \t %v \t %v", nodeWidth, nodeCenter)
		star := GetStar(starID)                                   // get the actual star
		starQuadrant := quadrant(star, nodeID)                    // find out in which quadrant it belongs
		quadrantNodeID := getQuadrantNodeID(nodeID, starQuadrant) // get the if of that quadrant
		insertIntoTree(starID, quadrantNodeID)                    // insert the star into that quadrant
	}
}

// directInsert inserts the star with the given ID into the given node inside of the given database
func directInsert(starID int64, nodeID int64) {
	// build the query
	query := fmt.Sprintf("UPDATE nodes SET star_id=%d WHERE node_id=%d", starID, nodeID)

	// Execute the query
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("[ E ] directInsert query: %v\n\t\t\t query: %s\n", err, query)
	}
}

// InsertList inserts all the stars in the given .csv into the stars and nodes table
func InsertList(database *sql.DB, filename string) {
	db = database
	// open the file
	content, readErr := ioutil.ReadFile(filename)
	if readErr != nil {
		panic(readErr)
	}

	in := string(content)
	reader := csv.NewReader(strings.NewReader(in))

	// insert all the stars into the db
	for {
		record, err := reader.Read()
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			log.Println("insertListErr")
			panic(err)
		}

		x, _ := strconv.ParseFloat(record[0], 64)
		y, _ := strconv.ParseFloat(record[1], 64)

		star := structs.Star2D{
			C: structs.Vec2{
				X: x / 100000,
				Y: y / 100000,
			},
			V: structs.Vec2{
				X: 0,
				Y: 0,
			},
			M: 1000,
		}

		fmt.Printf("Inserting (%f, %f)\n", star.C.X, star.C.Y)
		InsertStar(db, star, 1)
	}
}
