package main

import (
	"log"
	"os"

	"github.com/cbarraford/cryptocades-backend/api"
	"github.com/cbarraford/cryptocades-backend/manager"
	"github.com/cbarraford/cryptocades-backend/store"
	"github.com/cbarraford/cryptocades-backend/store/context"
)

// TODO: need a mechanism to shutdown the service for emergencies

func main() {
	var err error

	err = context.MigrateDB(os.Getenv("DATABASE_URL"), "file://./migrations")
	if err != nil {
		log.Fatal(err)
	}

	db, err := store.GetDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(116)
	db.SetMaxIdleConns(5)

	red, err := store.GetRedis(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	cstore := store.GetStore(db, red)

	manager.Start(cstore)

	r := api.GetAPIService(cstore)
	r.Run()
}
