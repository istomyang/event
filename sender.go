package event

import "sync"

const defaultSenderKey = "event_default_send_receiver"

func createDefaultSender() Sender {
	initDefaultMemoryCh()
	return NewMemorySender(defaultMemoryCh)
}

var defaultMemoryCh chan Message
var mut sync.Mutex

func initDefaultMemoryCh() {
	mut.Lock()
	defer mut.Unlock()
	if defaultMemoryCh == nil {
		defaultMemoryCh = make(chan Message)
	}
}
