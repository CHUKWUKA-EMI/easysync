package services

import (
	"fmt"
	"time"
)

// EventType ...
type EventType string

// Event ...
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

// EventListener ..
type EventListener chan<- Event

// EventEmitter ...
type EventEmitter struct {
	Listeners map[EventType]EventListener
}

// EventBus ...
var EventBus *EventEmitter

// EventTypes ...
const (
	SendEmailConfirmationCode EventType = "SendEmailConfirmationCode"
	ConfirmEmail              EventType = "ConfirmEmail"
)

// AddListener ...
func (evEmt *EventEmitter) AddListener(eventType EventType, listener EventListener) {
	fmt.Println("LISTENER ADDED!", listener)
	evEmt.Listeners[eventType] = listener
}

// Emit ...
func (evEmt *EventEmitter) Emit(event Event) {
	fmt.Println("EVENT EMITTED", event)
	listener, ok := evEmt.Listeners[event.Type]
	fmt.Println("LISTENER=", listener, "FOUND=", ok)
	listener <- event
}

// NewEventEmitter ...
func NewEventEmitter() *EventEmitter {
	evem := &EventEmitter{Listeners: make(map[EventType]EventListener)}
	EventBus = evem
	return evem
}
