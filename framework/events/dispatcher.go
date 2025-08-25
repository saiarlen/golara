package events

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

// Event represents an event
type Event interface {
	GetName() string
	GetPayload() map[string]interface{}
}

// BaseEvent provides basic event functionality
type BaseEvent struct {
	Name    string                 `json:"name"`
	Payload map[string]interface{} `json:"payload"`
}

func (e *BaseEvent) GetName() string {
	return e.Name
}

func (e *BaseEvent) GetPayload() map[string]interface{} {
	return e.Payload
}

// Listener represents an event listener
type Listener interface {
	Handle(event Event) error
}

// ListenerFunc is a function that implements Listener
type ListenerFunc func(event Event) error

func (f ListenerFunc) Handle(event Event) error {
	return f(event)
}

// EventDispatcher manages events and listeners
type EventDispatcher struct {
	listeners map[string][]Listener
	mutex     sync.RWMutex
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners: make(map[string][]Listener),
	}
}

// Listen registers a listener for an event
func (ed *EventDispatcher) Listen(eventName string, listener Listener) {
	ed.mutex.Lock()
	defer ed.mutex.Unlock()
	
	ed.listeners[eventName] = append(ed.listeners[eventName], listener)
}

// ListenFunc registers a function as a listener for an event
func (ed *EventDispatcher) ListenFunc(eventName string, handler func(event Event) error) {
	ed.Listen(eventName, ListenerFunc(handler))
}

// Dispatch dispatches an event to all listeners
func (ed *EventDispatcher) Dispatch(event Event) error {
	ed.mutex.RLock()
	listeners := ed.listeners[event.GetName()]
	ed.mutex.RUnlock()
	
	var errors []error
	
	for _, listener := range listeners {
		if err := ed.handleListener(listener, event); err != nil {
			errors = append(errors, err)
			log.Printf("Event listener error for %s: %v", event.GetName(), err)
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("event dispatch errors: %v", errors)
	}
	
	return nil
}

// DispatchAsync dispatches an event asynchronously
func (ed *EventDispatcher) DispatchAsync(event Event) {
	go func() {
		if err := ed.Dispatch(event); err != nil {
			log.Printf("Async event dispatch error: %v", err)
		}
	}()
}

func (ed *EventDispatcher) handleListener(listener Listener, event Event) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Event listener panic for %s: %v", event.GetName(), r)
		}
	}()
	
	return listener.Handle(event)
}

// RemoveListener removes a specific listener
func (ed *EventDispatcher) RemoveListener(eventName string, listener Listener) {
	ed.mutex.Lock()
	defer ed.mutex.Unlock()
	
	listeners := ed.listeners[eventName]
	for i, l := range listeners {
		if reflect.DeepEqual(l, listener) {
			ed.listeners[eventName] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

// RemoveAllListeners removes all listeners for an event
func (ed *EventDispatcher) RemoveAllListeners(eventName string) {
	ed.mutex.Lock()
	defer ed.mutex.Unlock()
	
	delete(ed.listeners, eventName)
}

// GetListeners returns all listeners for an event
func (ed *EventDispatcher) GetListeners(eventName string) []Listener {
	ed.mutex.RLock()
	defer ed.mutex.RUnlock()
	
	return ed.listeners[eventName]
}

// HasListeners checks if an event has listeners
func (ed *EventDispatcher) HasListeners(eventName string) bool {
	ed.mutex.RLock()
	defer ed.mutex.RUnlock()
	
	return len(ed.listeners[eventName]) > 0
}

// Common event types

// UserRegisteredEvent represents user registration
type UserRegisteredEvent struct {
	BaseEvent
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func NewUserRegisteredEvent(userID, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: BaseEvent{
			Name: "user.registered",
			Payload: map[string]interface{}{
				"user_id": userID,
				"email":   email,
			},
		},
		UserID: userID,
		Email:  email,
	}
}

// UserLoginEvent represents user login
type UserLoginEvent struct {
	BaseEvent
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
}

func NewUserLoginEvent(userID, ip string) *UserLoginEvent {
	return &UserLoginEvent{
		BaseEvent: BaseEvent{
			Name: "user.login",
			Payload: map[string]interface{}{
				"user_id": userID,
				"ip":      ip,
			},
		},
		UserID: userID,
		IP:     ip,
	}
}

// ModelCreatedEvent represents model creation
type ModelCreatedEvent struct {
	BaseEvent
	Model interface{} `json:"model"`
}

func NewModelCreatedEvent(model interface{}) *ModelCreatedEvent {
	return &ModelCreatedEvent{
		BaseEvent: BaseEvent{
			Name: "model.created",
			Payload: map[string]interface{}{
				"model": model,
			},
		},
		Model: model,
	}
}

// ModelUpdatedEvent represents model update
type ModelUpdatedEvent struct {
	BaseEvent
	Model interface{} `json:"model"`
}

func NewModelUpdatedEvent(model interface{}) *ModelUpdatedEvent {
	return &ModelUpdatedEvent{
		BaseEvent: BaseEvent{
			Name: "model.updated",
			Payload: map[string]interface{}{
				"model": model,
			},
		},
		Model: model,
	}
}

// ModelDeletedEvent represents model deletion
type ModelDeletedEvent struct {
	BaseEvent
	Model interface{} `json:"model"`
}

func NewModelDeletedEvent(model interface{}) *ModelDeletedEvent {
	return &ModelDeletedEvent{
		BaseEvent: BaseEvent{
			Name: "model.deleted",
			Payload: map[string]interface{}{
				"model": model,
			},
		},
		Model: model,
	}
}