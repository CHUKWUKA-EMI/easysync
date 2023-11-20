package main

import (
	"os"

	services "github.com/chukwuka-emi/easysync/Services"
)

func registerServices() {
	services.InitializeRedisConnection(os.Getenv("REDIS_URL"))

	// create event emitter/bus
	_ = services.NewEventEmitter()
}
