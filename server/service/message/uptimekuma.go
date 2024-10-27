package message

import (
	"log"
	"regexp"
	"sentinel/server/model"
	"strings"
)

type uptimeKuma struct {
	msgChan chan *model.Message
}

func (s *uptimeKuma) process() {
	for msg := range s.msgChan {
		uptimeMsg := s.ParseLogEntry(msg)
		if !s.isNormal(uptimeMsg) {
			log.Printf("Service Is Exception:", uptimeMsg.AppName)
		}
	}

}

// LogEntry 结构体表示提取的日志信息
type LogEntry struct {
	AppName string
	Status  string
	Msg     string
}

// ParseLogEntry 解析日志字符串并返回 LogEntry
func (s *uptimeKuma) ParseLogEntry(msg *model.Message) *LogEntry {
	// 更新的正则表达式，仅提取 Up 或 Down
	re := regexp.MustCompile(`^\[(.*?)\] \[\D*(Up|Down)\] (.*)$`)

	matches := re.FindStringSubmatch(msg.Message)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	return &LogEntry{
		AppName: matches[1],
		Status:  matches[2],
		Msg:     matches[3],
	}
}

func (s *uptimeKuma) isNormal(msg *LogEntry) bool {
	return strings.ToUpper(msg.Status) == "UP"
}