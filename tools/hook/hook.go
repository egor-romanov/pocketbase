package hook

import (
	"errors"
	"sync"
)

var StopPropagation = errors.New("Event hook propagation stopped")

// Handler defines a hook handler function.
type Handler[T any] func(data T) error

// Hook defines a concurrent safe structure for handling event hooks
// (aka. callbacks propagation).
type Hook[T any] struct {
	mux      sync.RWMutex
	handlers []Handler[T]
}

// Add registers a new handler to the hook.
func (h *Hook[T]) Add(fn Handler[T]) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.handlers = append(h.handlers, fn)
}

// Reset removes all registered handlers.
func (h *Hook[T]) Reset() {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.handlers = nil
}

// Trigger executes all registered hook handlers one by one
// with the specified `data` as an argument.
//
// Optionally, this method allows also to register additional one off
// handlers that will be temporary appended to the handlers queue.
//
// The execution stops when:
// - hook.StopPropagation is returned in one of the handlers
// - any non-nil error is returned in one of the handlers
func (h *Hook[T]) Trigger(data T, oneOffHandlers ...Handler[T]) error {
	h.mux.Lock()
	handlers := append(h.handlers, oneOffHandlers...)
	h.mux.Unlock() // unlock is not deferred to avoid deadlocks when Trigger is called recursive in the handlers

	for _, fn := range handlers {
		err := fn(data)
		if err == nil {
			continue
		}

		if errors.Is(err, StopPropagation) {
			return nil
		}

		return err
	}

	return nil
}
