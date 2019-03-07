package backend

import (
	"database/sql"
	"log"
)

// InitStarsTable initialises the stars table
func InitStarsTable(db *sql.DB) {
	log.Println("Preparing the query")
	var query = `CREATE TABLE IF NOT EXISTS public.stars
(
  star_id bigserial,
  x numeric,
  y numeric,
  vx numeric,
  vy numeric,
  m numeric,
  PRIMARY KEY (star_id)
) WITH (
  OIDS = FALSE
);

ALTER TABLE public.stars
  OWNER to postgres;`
	log.Println("Executing the query")
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("[ E ] InitNodesTable query: %v \n\t\t\tquery: %s\n", err, query)
	}
	log.Println("DONE")
}

// InitNodesTable initialises the nodes table
func InitNodesTable(db *sql.DB) {
	log.Println("creating the query")
	var query = `CREATE TABLE IF NOT EXISTS public.nodes
(
  node_id bigserial NOT NULL,
  box_width numeric,
  total_mass numeric,
  depth integer,
  star_id bigint,
  root_id bigint,
  isleaf boolean,
  box_center numeric[],
  center_of_mass numeric[],
  subnodes bigint[],
  PRIMARY KEY (node_id)
) WITH (
  OIDS = FALSE
);

ALTER TABLE public.nodes
  OWNER to postgres;`
	log.Println("done creating the query")
	log.Println("executing the query")
	_, err := db.Exec(query)
	log.Println("done executing the query")
	if err != nil {
		log.Fatalf("[ E ] InitNodesTable query: %v \n\t\t\tquery: %s\n", err, query)
	}
}
