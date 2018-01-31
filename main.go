package main

import (
	"fmt"
	"log"
	"os"

	newrelic "github.com/newrelic/go-agent"

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

	agentName := fmt.Sprintf("Cryptocades-%s", os.Getenv("ENVIRONMENT"))
	key := os.Getenv("NEW_RELIC_LICENSE_KEY")
	agentConfig := newrelic.NewConfig(agentName, key)
	// agentConfig.Logger = newrelic.NewDebugLogger(os.Stdout)
	// if we don't have a license key for new relic, disable it.
	if len(key) == 0 {
		agentConfig.Enabled = false
	}
	agent, err := newrelic.NewApplication(agentConfig)
	if err != nil {
		log.Fatal(err)
	}

	r := api.GetAPIService(cstore, agent)
	r.Run()
}
