package message

import "sentinel/server/model"

type processor interface {
	process()
	push(msg *model.Message)
}

var processorMap map[model.ProcessorAPPID]processor

func RegisterProcessor() {
	processorMap = make(map[model.ProcessorAPPID]processor, 16)
	uk := &uptimeKuma{msgChan: make(chan *model.Message, 100)}
	processorMap[1] = uk
	go uk.process()
}

func SendMsg(msg *model.Message) {
	processorMap[msg.APPID].push(msg)
}
