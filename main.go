package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ruyoutor/simplebank/api"
	db "github.com/ruyoutor/simplebank/db/sqlc"
	"github.com/ruyoutor/simplebank/util"
)

var config util.Config

func init() {
	var err error
	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
}

func main() {

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
