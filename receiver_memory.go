package event

type memoryReceiver struct {
	receiver <-chan Message
}

func NewMemoryReceiver(ch <-chan Message) Receiver {
	return &memoryReceiver{receiver: ch}
}

func (f *memoryReceiver) Receive() <-chan Message {
	return f.receiver
}

func (f *memoryReceiver) Run() error {
	return nil
}

func (f *memoryReceiver) MustRun() {
	return
}

func (f *memoryReceiver) Close() error {
	return nil
}

var _ Receiver = &memoryReceiver{}
