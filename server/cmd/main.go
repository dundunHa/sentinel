package main

import (
	"log"
	"net/http"
	"sentinel/server/api"
	"sentinel/server/scheduler"
	"sentinel/server/service"
	"sentinel/server/service/message"
	"sentinel/server/storage"
)

func main() {

	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	svc := service.NewService()
	sched := scheduler.NewScheduler(svc)
	message.RegisterProcessor()
	if err := service.RegisterDockerClient(); err != nil {
		log.Fatalf("Failed to register docker client: %v", err)
	}

	go sched.Start()

	r := api.NewAPI(svc).Routes()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
