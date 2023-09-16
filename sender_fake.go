package event

import (
	"fmt"
	"sync/atomic"
)

type fakeSender struct {
	sender chan<- Message
	closed atomic.Bool
}

func newFakeSender(ch chan<- Message) Sender {
	return &fakeSender{sender: ch}
}

func (f *fakeSender) Name() string {
	return "FakeReceiver"
}

func (f *fakeSender) Send(message Message) {
	fmt.Printf("%s: Send, got %v.\n", f.Name(), message)
	go func() {
		if f.closed.Load() {
			return
		}
		f.sender <- message
	}()
}

func (f *fakeSender) SendSync(message Message) error {
	fmt.Printf("%s: SendSync, got %v.\n", f.Name(), message)
	if f.closed.Load() {
		return nil
	}
	f.sender <- message
	return nil
}

func (f *fakeSender) Run() error {
	fmt.Printf("%s: Run().\n", f.Name())
	return nil
}

func (f *fakeSender) MustRun() {
	fmt.Printf("%s: MustRun().\n", f.Name())
	return
}

func (f *fakeSender) Close() error {
	fmt.Printf("%s: Close().\n", f.Name())
	f.closed.Store(true)
	return nil
}

var _ Sender = &fakeSender{}
