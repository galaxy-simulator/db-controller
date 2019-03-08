package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"git.darknebu.la/GalaxySimulator/structs"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.darknebu.la/GalaxySimulator/db-controller/backend"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func requestInfo(r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.Host, r.URL.Path)
}

func errHandler(name string, err error) {
	if err != nil {
		log.Fatalf("Error: %s: %v", name, err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)
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
	requestInfo(r)

	width, _ := strconv.ParseFloat(r.Form.Get("w"), 64)
	backend.NewTree(db, width)
}

func deleteStarsHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	backend.DeleteAllNodes(db)
}

func deleteNodesHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	backend.DeleteAllStars(db)
}

func starslistGoHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	listOfStarsGo := backend.GetListOfStarsGo(db)
	_, _ = fmt.Fprintf(w, "%v", listOfStarsGo)
}

func starslistCsvHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	listOfStarsCsv := backend.GetListOfStarsCsv(db)
	for _, element := range listOfStarsCsv {
		_, _ = fmt.Fprintf(w, "%s\n", element)
	}
}

func updateTotalMassHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	backend.UpdateTotalMass(db, index)
}

func updateCenterOfMassHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	backend.UpdateCenterOfMass(db, index)
}

func genForestTreeHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	parseFormErr := r.ParseForm()
	errHandler("genForestTree Parse Form", parseFormErr)

	// get the index of the tree that should be shown
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	log.Printf("Generating the forest representation of the tree with the index %d", index)

	// generate the forest representation and write it back to the user
	forestTree := backend.GenForestTree(db, index)
	_, _ = fmt.Fprintf(w, "%s\n", forestTree)
}

func createNodesTableHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	backend.InitNodesTable(db)
}

func createStarsTableHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	backend.InitStarsTable(db)
}

func insertStarHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	// parse the http post form
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	log.Printf("x: %s", r.PostFormValue("x"))
	log.Printf("y: %s", r.PostFormValue("y"))
	log.Printf("vx: %s", r.PostFormValue("vx"))
	log.Printf("vy: %s", r.PostFormValue("vy"))
	log.Printf("m: %s", r.PostFormValue("m"))

	// parse the star parameters
	x, xParseErr := strconv.ParseFloat(r.PostFormValue("x"), 64)
	errHandler("parse x", xParseErr)
	y, yParseErr := strconv.ParseFloat(r.PostFormValue("y"), 64)
	errHandler("parse y", yParseErr)
	vx, vxParseErr := strconv.ParseFloat(r.PostFormValue("vx"), 64)
	errHandler("parse vx", vxParseErr)
	vy, vyParseErr := strconv.ParseFloat(r.PostFormValue("vy"), 64)
	errHandler("parse vy", vyParseErr)
	m, mParseErr := strconv.ParseFloat(r.PostFormValue("m"), 64)
	errHandler("parse m", mParseErr)

	// build the star
	var star = structs.Star2D{
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

	// parse the tree index
	index, indexParseErr := strconv.ParseInt(r.Form.Get("index"), 10, 64)
	errHandler("parse index", indexParseErr)

	// Insert the star into the tree
	backend.InsertStar(db, star, index)
}

func insertStarListHandler(w http.ResponseWriter, r *http.Request) {
	requestInfo(r)

	// parse the http post form
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// get the filename of the file to insert
	filename := r.PostFormValue("filename")

	// open the csv file and parse it
	csvFile, _ := os.Open(fmt.Sprintf("%s.csv", filename))
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// iterate over all the lines
	for {

		// read the line
		line, error := reader.Read()

		// handler errors such as broken syntax and the end of the file
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		// parse the star parameters
		x, xParseErr := strconv.ParseFloat(line[0], 64)
		errHandler("parse x", xParseErr)
		y, xParseErr := strconv.ParseFloat(line[1], 64)
		errHandler("parse y", xParseErr)

		// define a star
		star := structs.Star2D{
			C: structs.Vec2{
				X: x,
				Y: y,
			},
			V: structs.Vec2{
				X: 0,
				Y: 0,
			},
			M: 1000,
		}

		// insert the star into the database
		backend.InsertStar(db, star, 1)
	}
}

func main() {
	// get the port on which the service should be hosted and the url of the database
	var port string
	flag.StringVar(&port, "port", "8080", "port used to host the service")
	flag.Parse()
	log.Println("[ ] Done loading the flags")

	// get the data that should be used to connect to the database
	var DBURL = os.Getenv("DBURL")
	var DBUSER = os.Getenv("DBUSER")
	var DBPASSWD = os.Getenv("DBPASSWD")
	var DBPORT, _ = strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	var DBPROJECTNAME = os.Getenv("DBPROJECTNAME")

	log.Printf("DBURL: %s", DBURL)
	log.Printf("DBUSER: %s", DBUSER)
	log.Printf("DBPASSWD: %s", DBPASSWD)
	log.Printf("DBPORT: %d", DBPORT)
	log.Printf("DBPROJECTNAME: %s", DBPROJECTNAME)
	log.Printf("frontend port: %s", port)

	// connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBURL, DBPORT, DBUSER, DBPASSWD, DBPROJECTNAME)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error: The data source arguments are not valid: %s", err)
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
	router.HandleFunc("/insertStar", insertStarHandler).Methods("POST")
	router.HandleFunc("/insertStarList", insertStarListHandler).Methods("POST")

	log.Printf("[ ] Starting the API on localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
