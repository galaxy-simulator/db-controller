package backend

import (
	"database/sql"
	"log"
)

// InitStarsTable initialises the stars table
func InitStarsTable(db *sql.DB) {
	query := `CREATE TABLE public.stars
(
    star_id bigint NOT NULL DEFAULT nextval('stars_star_id_seq'::regclass),
    x numeric,
    y numeric,
    vx numeric,
    vy numeric,
    m numeric
)
`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("[ E ] InitNodesTable query: %v \n\t\t\tquery: %s\n", err, query)
	}
}

// InitNodesTable initialises the nodes table
func InitNodesTable(db *sql.DB) {
	query := `CREATE TABLE public.nodes
	(
		node_id bigint NOT NULL DEFAULT nextval('nodes_node_id_seq'::regclass),
	box_width numeric NOT NULL,
		total_mass numeric NOT NULL,
		depth integer,
		star_id bigint NOT NULL,
		root_id bigint NOT NULL,
		isleaf boolean,
		box_center numeric[] NOT NULL,
		center_of_mass numeric[] NOT NULL,
		subnodes bigint[] NOT NULL
	)
`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("[ E ] InitNodesTable query: %v \n\t\t\tquery: %s\n", err, query)
	}
}
