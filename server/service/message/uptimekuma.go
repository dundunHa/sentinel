package message

import (
	"log"
	"regexp"
	"sentinel/server/model"
	"strings"
	"sentinel/server/service/docker"
)

type uptimeKuma struct {
	msgChan chan *model.Message
}

func (s *uptimeKuma) process() {
	for msg := range s.msgChan {
		uptimeMsg := s.ParseLogEntry(msg)
		if uptimeMsg == nil {
			continue
		}
		if !s.isNormal(uptimeMsg) {
			log.Println("Service Is Exception:", uptimeMsg.AppName)
			log.Println("will get docker client")
			dkservice :=docker.NewDockerService(docker.DockerClientMap[msg.APPID])
			log.Println("cli:",dkservice.Cli)
			dkservice.RestartContainer(uptimeMsg.AppName)
		} else {
			log.Println("Service Is Running:", uptimeMsg.AppName)
		}
	}

}

func (s *uptimeKuma) push(msg *model.Message) {
	s.msgChan <- msg
}

type LogEntry struct {
	AppName     string
	Status      string
	Msg         string
	GotifyAppID model.ProcessorAPPID
}

func (s *uptimeKuma) ParseLogEntry(msg *model.Message) *LogEntry {
	re := regexp.MustCompile(`^\[(.*?)\] \[\D*(Up|Down)\] (.*)$`)
	matches := re.FindStringSubmatch(msg.Message)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	return &LogEntry{
		AppName:     matches[1],
		Status:      matches[2],
		Msg:         matches[3],
		GotifyAppID: msg.APPID,
	}
}

func (s *uptimeKuma) isNormal(msg *LogEntry) bool {
	return strings.ToUpper(msg.Status) == "UP"
}
