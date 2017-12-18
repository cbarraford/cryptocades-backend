package main

import (
	"log"
	"os"

	"github.com/CBarraford/lotto/api"
	"github.com/CBarraford/lotto/manager"
	"github.com/CBarraford/lotto/store"
	"github.com/CBarraford/lotto/store/context"
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
	cstore := store.GetStore(db)

	manager.Start(cstore)

	r := api.GetAPIService(cstore)
	r.Run()
}
