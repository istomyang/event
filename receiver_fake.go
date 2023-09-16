package event

import (
	"fmt"
)

type fakeReceiver struct {
	receiver <-chan Message
}

func newFakeReceiver(ch <-chan Message) Receiver {
	return &fakeReceiver{receiver: ch}
}

func (f *fakeReceiver) Name() string {
	return "FakeReceiver"
}

func (f *fakeReceiver) Receive() <-chan Message {
	return f.receiver
}

func (f *fakeReceiver) Run() error {
	fmt.Printf("%s: Run().\n", f.Name())
	return nil
}

func (f *fakeReceiver) MustRun() {
	fmt.Printf("%s: MustRun().\n", f.Name())
	return
}

func (f *fakeReceiver) Close() error {
	fmt.Printf("%s: Close().\n", f.Name())
	return nil
}

var _ Receiver = &fakeReceiver{}
