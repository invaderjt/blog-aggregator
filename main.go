package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	state := initialize()
	db, err := sql.Open("postgres", state.Cfg.DBURL)
	if err != nil {
		log.Fatalln(err)
	}

	updateState(state, db)

}
