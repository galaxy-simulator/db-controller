package backend

import (
	"database/sql"
	"git.darknebu.la/GalaxySimulator/structs"
	"log"
	"math"
)

// CalcAllForces calculates all the forces acting on the given star.
// The theta value it receives is used by the Barnes-Hut algorithm to determine what
// stars to include into the calculations
func CalcAllForces(database *sql.DB, star structs.Star2D, galaxyIndex int64, theta float64) structs.Vec2 {
	db = database

	// calculate all the forces and add them to the list of all forces
	// this is done recursively
	// first of all, get the root id
	log.Println("[db_actions] Getting the root ID")
	rootID := getRootNodeID(galaxyIndex)
	log.Println("[db_actions] Done getting the root ID")

	log.Printf("[db_actions] Calculating the forces acting on the star %v", star)
	force := CalcAllForcesNode(star, rootID, theta)
	log.Printf("[db_actions] Done calculating the forces acting on the star %v", star)
	log.Printf("[db_actions] Force: %v", force)

	return force
}

// calcAllForces nodes calculates the forces in between a sta	log.Printf("Calculating the forces acting on the star %v", star)r and a node and returns the overall force
// TODO: implement the calcForce(star, centerOfMass) {...} function
// TODO: implement the getSubtreeIDs(nodeID) []int64 {...} function
func CalcAllForcesNode(star structs.Star2D, nodeID int64, theta float64) structs.Vec2 {
	log.Println("---------------------------------------")
	log.Printf("NodeID: %d \t star: %v \t theta: %f \t nodeboxwidth: %f", nodeID, star, theta, getBoxWidth(nodeID))
	var forceX float64
	var forceY float64
	var localTheta float64

	nodeWidth := getBoxWidth(nodeID)

	if nodeID != 0 {
		log.Println("[theta] Calculating localtheta(star, node)")
		log.Printf("[theta] node with: %f", nodeWidth)
		localTheta = calcTheta(star, nodeID)
		log.Printf("[theta] Done calculating localtheta: %v", localTheta)
	}

	// recurse deeper into the tree
	if localTheta < theta {
		log.Println("[   ] localtheta < theta")

	} else {
		log.Println("[   ] localtheta > theta")

		log.Printf("[   ] Iterating over subtrees")
		var subtreeIDs [4]int64
		subtreeIDs = getSubtreeIDs(nodeID)
		for i, subtreeID := range subtreeIDs {
			log.Printf("Subtree: %d\t ID: %d", i, subtreeID)

			if subtreeID != 0 {
				subtreeStarId := getStarID(subtreeID)
				if subtreeStarId != 0 {
					var localStar = GetStar(subtreeStarId)
					log.Printf("subtree %d star: %v", i, localStar)
					if localStar != star {
						log.Println("Not even the original star, calculating forces...")
						var force = calcForce(localStar, star)
						forceX += force.X
						forceY += force.Y
					}
				}
				var force = CalcAllForcesNode(star, subtreeID, theta)
				log.Printf("force: %v", force)
				forceX += force.X
				forceY += force.Y
			}
		}

	}

	log.Println("---------------------------------------")
	return structs.Vec2{forceX, forceY}
}

// calcTheta calculates the theat for a given star and a node
func calcTheta(star structs.Star2D, nodeID int64) float64 {
	d := getBoxWidth(nodeID)
	r := distance(star, nodeID)
	theta := d / r
	return theta
}

// calculate the distance in between the star and the node with the given ID
func distance(star structs.Star2D, nodeID int64) float64 {
	var starX float64 = star.C.X
	var starY float64 = star.C.Y
	var node structs.Vec2 = getNodeCenterOfMass(nodeID)
	var nodeX float64 = node.X
	var nodeY float64 = node.Y

	var tmpX = math.Pow(starX-nodeX, 2)
	var tmpY = math.Pow(starY-nodeY, 2)

	var distance float64 = math.Sqrt(tmpX + tmpY)
	return distance
}

// calcForce calculates the force the star s1 is acting on s2.
// The force acting is returned in Newtons.
func calcForce(s1 structs.Star2D, s2 structs.Star2D) structs.Vec2 {
	log.Println("+++++++++++++++++++++++++")
	log.Printf("s1: %v", s1)
	log.Printf("s2: %v", s2)
	G := 6.6726 * math.Pow(10, -11)

	// calculate the force acting
	var combinedMass float64 = s1.M * s2.M
	var distance float64 = math.Sqrt(math.Pow(math.Abs(s1.C.X-s2.C.X), 2) + math.Pow(math.Abs(s1.C.Y-s2.C.Y), 2))
	log.Printf("combined mass: %f", combinedMass)
	log.Printf("distance: %f", distance)

	var scalar float64 = G * ((combinedMass) / math.Pow(distance, 2))
	log.Printf("scalar: %f", scalar)

	// define a unit vector pointing from s1 to s2
	var vector structs.Vec2 = structs.Vec2{s2.C.X - s1.C.X, s2.C.Y - s1.C.Y}
	var UnitVector structs.Vec2 = structs.Vec2{vector.X / distance, vector.Y / distance}

	// multiply the vector with the force to get a vector representing the force acting
	var force structs.Vec2 = UnitVector.Multiply(scalar)
	log.Println("+++++++++++++++++++++++++")

	// return the force exerted on s1 by s2
	return force
}
