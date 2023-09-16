package event

const defaultReceiverKey = "event_default_send_receiver"

func createDefaultReceiver() Receiver {
	initDefaultMemoryCh()
	return NewMemoryReceiver(defaultMemoryCh)
}
