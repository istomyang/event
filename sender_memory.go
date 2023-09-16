package event

import (
	"sync/atomic"
)

type memorySender struct {
	sender chan<- Message
	closed atomic.Bool
}

func NewMemorySender(ch chan<- Message) Sender {
	return &memorySender{sender: ch}
}

func (f *memorySender) Send(message Message) {
	go func() {
		if f.closed.Load() {
			return
		}
		f.sender <- message
	}()
}

func (f *memorySender) SendSync(message Message) error {
	if f.closed.Load() {
		return nil
	}
	f.sender <- message
	return nil
}

func (f *memorySender) Run() error {
	return nil
}

func (f *memorySender) MustRun() {
	return
}

func (f *memorySender) Close() error {
	f.closed.Store(true)
	close(f.sender)
	return nil
}

var _ Sender = &memorySender{}
