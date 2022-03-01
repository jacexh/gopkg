package mediator

import (
	"context"
	"sync/atomic"
)

type (
	// EventType 事件类型描述
	EventType string

	// EventHandler 事件处理函数
	EventHandler func(context.Context, Event)

	// Event 事件接口
	Event interface {
		Type() EventType
	}

	eventCollection struct {
		events []Event
		raised int32
	}

	EventCollection interface {
		Add(Event)
		Raise(context.Context, Mediator)
	}
)

func NewEventCollection() EventCollection {
	return &eventCollection{events: make([]Event, 0)}
}

func (es *eventCollection) Add(ev Event) {
	if atomic.LoadInt32(&es.raised) == 0 {
		es.events = append(es.events, ev)
	}
}

func (es *eventCollection) Raise(ctx context.Context, m Mediator) {
	if atomic.CompareAndSwapInt32(&es.raised, 0, 1) {
		for _, event := range es.events {
			m.Dispatch(ctx, event)
		}
	}
}
