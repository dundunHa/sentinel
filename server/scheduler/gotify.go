package scheduler

import (
	"log"
	"sentinel/server/service"
	"time"
)

type Scheduler struct {
	svc *service.Service
}

func NewScheduler(svc *service.Service) *Scheduler {
	return &Scheduler{svc: svc}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Checking for new messages...")
			if err := s.svc.HandleMessages(); err != nil {
				log.Printf("Error handling messages: %v", err)
			}
		}
	}
}
