package ogo

import "fmt"

type HandlerLogger struct {
	path string
	ch   chan string
}

func NewHandlerLogger(path string) *HandlerLogger {
	ch := make(chan string)
	go func() {
		fmt.Println(<-ch)
	}()
	return &HandlerLogger{
		path: path,
		ch:   ch,
	}
}

func (h *HandlerLogger) log(format string, a ...any) {
	h.ch <- fmt.Sprintf("OGO: handler: %v\n%v", h.path, fmt.Sprintf(format, a...))
}

// logs and panics
func (h *HandlerLogger) logAndPanic(format string, a ...any) {
	msg := fmt.Sprintf("OGO: handler: %v\n%v", h.path, fmt.Sprintf(format, a...))
	h.ch <- msg
	panic(msg)
}
