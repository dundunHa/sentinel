package message

import "sentinel/server/model"

type processor interface {
	process()
}

var processorMap map[string]processor

func RegisterProcessor() {
	processorMap = make(map[string]processor, 16)
	uk := &uptimeKuma{msgChan: make(chan *model.Message, 100)}

	processorMap["uptimekuma"] = uk
}

func SendMsg(msg *model.Message) {

}
