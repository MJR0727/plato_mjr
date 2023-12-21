package source

import (
	"fmt"
	"hello/plato_mjr/common/discovery"
)

type EventType string

type Event struct {
	Type          EventType
	IP            string
	Port          string
	MessageBytes  float64 // node current status
	ConnectionNum float64
}

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

var eventChan chan *Event

func EventChan() <-chan *Event {
	return eventChan
}

func (event *Event) Key() string {
	return fmt.Sprintf("%s:%s", event.IP, event.Port)
}

func NewEvent(epi *discovery.EndpointInfo) *Event {
	var messageByte, connectionNum float64
	if data, ok := epi.MetaData["message_byte"]; ok {
		messageByte = data.(float64) // 如果出错，此处应该panic 暴露错误
	}
	if data, ok := epi.MetaData["connet_num"]; ok {
		connectionNum = data.(float64) // 如果出错，此处应该panic 暴露错误
	}
	var event = &Event{
		IP:            epi.IP,
		Port:          epi.Port,
		MessageBytes:  messageByte,
		ConnectionNum: connectionNum,
		Type:          AddNodeEvent,
	}
	return event
}
