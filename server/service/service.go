package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sentinel/server/config"
	"sentinel/server/dao"
	"sentinel/server/model"
	msgProcessor "sentinel/server/service/message"
)

type Service struct {
	client        *GotifyClient
	dockerService *DockerService
}

func NewService() *Service {
	cfg := config.LoadConfig()
	client := NewClient(cfg.GotifyURL, cfg.GotifyToken)
	dkService, err := NewDockerService()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return &Service{
		client:        client,
		dockerService: dkService,
	}
}

func (s *Service) HandleMessages() error {
	lastProcessedID, err := dao.GetLastProcessedID(s.client.appID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	messages, err := s.GetMessages()
	if err != nil {
		return err
	}

	var maxID = lastProcessedID
	for _, message := range messages {
		if message.ID <= lastProcessedID {
			continue
		}
		msgProcessor.SendMsg(&message)
		if message.ID > maxID {
			maxID = message.ID
		}
	}

	if maxID > lastProcessedID {
		if err := dao.UpdateLastProcessedID(s.client.appID, maxID); err != nil {
			log.Printf("Failed to update last processed ID: %v", err)
		}
	}
	return nil
}

type GotifyClient struct {
	url   string
	token string
	appID int
}

func NewClient(url, token string) *GotifyClient {
	return &GotifyClient{url: url, token: token}
}

func (c *GotifyClient) SendMessage(title, message string) error {
	payload := map[string]string{"title": title, "message": message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url+"/message", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.New("failed to send message")
	}
	return nil
}

func (s *Service) GetMessages() ([]model.Message, error) {
	url := fmt.Sprintf("%s/message", s.client.url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Gotify-Key", s.client.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get messages")
	}

	var gotifyRsp model.GotifyMsg
	if err := json.NewDecoder(resp.Body).Decode(&gotifyRsp); err != nil {
		return nil, err
	}
	return gotifyRsp.Messages, nil
}
