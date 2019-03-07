package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.darknebu.la/GalaxySimulator/db-controller/backend"
)

var (
	db *sql.DB
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexString := `<html><body><h1>Galaxy Simulator Database Frontend</h1>

		API:
	<h3> / (GET) </h3>

	<h3> /new (POST) </h3>
		Create a new Tree
	<br>
		Parameters:
	<ul>
	<li>
		w float64: width of the tree
	</li>
	</ul>

	<h3> /deleteStars (POST) </h3>
		Delete all stars from the stars Table
	<br>
		Parameters:
	<ul>
	<li>
		none
	</li>
	</ul>

	<h3> /deleteNodes (POST) </h3>
		Delete all nodes from the nodes Table
	<br>
		Parameters:
	<ul>
	<li>
		none
	</li>
	</ul>

	<h3> /starslist/go (GET) </h3>
		List all stars using go-array format
	<br>
		Parameters:
	<ul>
	<li>
		none
	</li>
	</ul>

	<h3> /starslist/csv (GET) </h3>
		List all stars as a csv
	<br>
		Parameters:
	<ul>
	<li>
		none
	</li>
	</ul>

	<h3> /updatetotalmass (POST) </h3>
		Update the total mass of all the nodes in the tree with the given index
	<br>
		Parameters:
	<ul>
	<li>
		index int: index of the tree
	</li>
	</ul>

	<h3> /updatecenterofmass (POST) </h3>
		Update the center of mass of all the nodes in the tree with the given index
	<br>
		Parameters:
	<ul>
	<li>
		index int: index of the tree
	</li>
	</ul>

	<h3> /genforesttree (GET) </h3>
		Generate the forest representation of the tree with the given index
	<br>
		Parameters:
	<ul>
	<li>
		index int: index of the tree
	</li>
	</ul>

	</body>
	</html>
`
	_, _ = fmt.Fprintf(w, "%s", indexString)
}
func newTreeHandler(w http.ResponseWriter, r *http.Request) {
	width, _ := strconv.ParseFloat(r.Form.Get("w"), 64)
	backend.NewTree(db, width)
}

func deleteStarsHandler(w http.ResponseWriter, r *http.Request) {
	backend.DeleteAllNodes(db)
}

func deleteNodesHandler(w http.ResponseWriter, r *http.Request) {
	backend.DeleteAllStars(db)
}

func starslistGoHandler(w http.ResponseWriter, r *http.Request) {
	backend.GetListOfStarsGo(db)
}

func starslistCsvHandler(w http.ResponseWriter, r *http.Request) {
	backend.GetListOfStarsCsv(db)
}

func updateTotalMassHandler(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	backend.UpdateTotalMass(db, index)
}

func updateCenterOfMassHandler(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	backend.UpdateCenterOfMass(db, index)
}

func genForestTreeHandler(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	backend.GenForestTree(db, index)
}

func createNodesTableHandler(w http.ResponseWriter, r *http.Request) {
	backend.InitNodesTable(db)
}

func createStarsTableHandler(w http.ResponseWriter, r *http.Request) {
	backend.InitStarsTable(db)
}

func main() {
	// get the port on which the service should be hosted and the url of the database
	var port string
	flag.StringVar(&port, "port", "8080", "port used to host the service")
	var dbURL string
	flag.StringVar(&dbURL, "DBURL", "postgres", "url of the database used")
	flag.Parse()

	log.Println("[ ] Done loading the flags")

	// get the data that should be used to connect to the database
	var DBUSER = os.Getenv("DBUSER")
	var DBPASSWD = os.Getenv("DBPASSWD")
	var DBPORT, _ = strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	var DBPROJECTNAME = os.Getenv("DBPROJECTNAME")

	log.Println("[ ] Done loading the envs")

	// connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbURL, DBPORT, DBUSER, DBPASSWD, DBPROJECTNAME)
	log.Printf("[ ] Done assembling the connString: %s", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		_ = fmt.Errorf("%s", err)
	}

	log.Println("[ ] Done Connecting to the DB")

	// ping the db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("[ ] Done Pinging the DB")

	// define a new mux router
	router := mux.NewRouter()

	// define the endpoints
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/newTree", newTreeHandler).Methods("POST")
	router.HandleFunc("/deleteStars", deleteStarsHandler).Methods("POST")
	router.HandleFunc("/deleteNodes", deleteNodesHandler).Methods("POST")
	router.HandleFunc("/starslist/go", starslistGoHandler).Methods("GET")
	router.HandleFunc("/starslist/csv", starslistCsvHandler).Methods("GET")
	router.HandleFunc("/updateTotalMass", updateTotalMassHandler).Methods("POST")
	router.HandleFunc("/updateCenterOfMass", updateCenterOfMassHandler).Methods("POST")
	router.HandleFunc("/genForestTree", genForestTreeHandler).Methods("GET")
	router.HandleFunc("/createNodesTable", createNodesTableHandler).Methods("POST")
	router.HandleFunc("/createStarsTable", createStarsTableHandler).Methods("POST")

	// start the http server on the port reached in via a flag
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
