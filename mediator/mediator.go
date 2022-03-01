package mediator

import "context"

type (
	Mediator interface {
		Dispatch(context.Context, Event)
		Subscribe(EventType, EventHandler)
	}

	inMemMediator struct {
		handlers           map[EventType][]EventHandler
		concurrent         chan struct{}
		orphanEventHandler func(Event)
	}

	Option func(*inMemMediator)
)

var (
	_ Mediator = (*inMemMediator)(nil)
)

func NewInMemMediator(c int, opts ...Option) Mediator {
	m := &inMemMediator{
		handlers:   make(map[EventType][]EventHandler),
		concurrent: make(chan struct{}, c),
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *inMemMediator) Subscribe(et EventType, hdl EventHandler) {
	if _, ok := m.handlers[et]; !ok {
		m.handlers[et] = make([]EventHandler, 0)
	}
	m.handlers[et] = append(m.handlers[et], hdl)
}

func (m *inMemMediator) Dispatch(ctx context.Context, ev Event) {
	if _, ok := m.handlers[ev.Type()]; !ok {
		if m.orphanEventHandler != nil {
			m.orphanEventHandler(ev)
			return
		}
		return
	}

	m.concurrent <- struct{}{}
	go func(ctx context.Context, ev Event, handlers ...EventHandler) { // 确保event的多个handler处理的顺序以及时效性
		defer func() {
			<-m.concurrent
		}()
		for _, handler := range handlers {
			handler(ctx, ev) // 在handler内部处理ctx.Done()
		}
	}(ctx, ev, m.handlers[ev.Type()]...)
}

func WithOrphanEventHandler(fn func(Event)) Option {
	return func(m *inMemMediator) {
		m.orphanEventHandler = fn
	}
}
