package event

import (
	"strings"
	"sync"
)

type subscribe struct {
	config                          *subscribeConfig
	fullPath                        string
	messageKeyHandlerMap            map[string][]SubscribeHandler
	receiverKeyMessageKeyHandlerMap map[string]map[string][]SubscribeHandler
	childrenMap                     map[string]*subscribe // HashMap
	mut                             sync.RWMutex
	fastCheck                       bool
}

type subscribeConfig struct {
}

var _ Subscribe = &subscribe{}

func newSubscribe(config subscribeConfig) subscribe {
	return subscribe{
		config:                          &config,
		messageKeyHandlerMap:            make(map[string][]SubscribeHandler),
		receiverKeyMessageKeyHandlerMap: make(map[string]map[string][]SubscribeHandler),
		childrenMap:                     make(map[string]*subscribe),
		fastCheck:                       true,
	}
}

func (s *subscribe) Handle(context *Context) {
	// Handler in this Layer.
	if !s.fastCheck || context.Message.Metadata.FullPath == s.FullPathString() {
		receiverKey := context.Message.Metadata.ChannelKey
		messageKey := context.Message.Metadata.Key
		var handlers []SubscribeHandler
		s.mut.RLock()
		if receiverKey == "" {
			handlers = s.messageKeyHandlerMap[messageKey]
		} else {
			handlers = s.receiverKeyMessageKeyHandlerMap[receiverKey][messageKey]
		}
		s.mut.RUnlock()
		for _, handler := range handlers {
			handler(context)
		}
		return
	}
	// handle in next node layer, here use hashmap.
	sub := s.childrenMap[context.Message.Metadata.FullPath]
	sub.Handle(context) // if sub is nil, not panic, because this call do THandle(sub *T, context).
	// Not found in root layer, return to discard.
}

func (s *subscribe) Group(path string) Subscribe {
	var pre = s.fullPath
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	fullPath := pre + "/" + path
	sub := &subscribe{
		config:                          s.config,
		fullPath:                        fullPath,
		messageKeyHandlerMap:            make(map[string][]SubscribeHandler),
		receiverKeyMessageKeyHandlerMap: make(map[string]map[string][]SubscribeHandler),
		childrenMap:                     s.childrenMap,
	}
	sub.childrenMap[fullPath] = sub
	return sub
}

func (s *subscribe) Subscribe(messageKey string, handler SubscribeHandler) {
	s.mut.Lock()
	s.messageKeyHandlerMap[messageKey] = append(s.messageKeyHandlerMap[messageKey], handler)
	s.mut.Unlock()
}

func (s *subscribe) SubscribeKeys(handler SubscribeHandler, messageKeys ...string) {
	for _, key := range messageKeys {
		s.Subscribe(key, handler)
	}
}

func (s *subscribe) UnSubscribe(messageKey string) {
	s.mut.Lock()
	delete(s.messageKeyHandlerMap, messageKey)
	s.mut.Unlock()
}

func (s *subscribe) UnSubscribeKeys(messageKeys ...string) {
	for _, key := range messageKeys {
		s.UnSubscribe(key)
	}
}

func (s *subscribe) SubscribeReceiver(receiverKey string, messageKey string, handler SubscribeHandler) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if _, has := s.receiverKeyMessageKeyHandlerMap[receiverKey]; !has {
		s.receiverKeyMessageKeyHandlerMap[receiverKey] = make(map[string][]SubscribeHandler)
	}
	m := s.receiverKeyMessageKeyHandlerMap[receiverKey]
	m[messageKey] = append(m[messageKey], handler)
}

func (s *subscribe) UnSubscribeReceiver(receiverKey string, messageKey string) {
	s.mut.Lock()
	m := s.receiverKeyMessageKeyHandlerMap[receiverKey]
	delete(m, messageKey)
	s.mut.Unlock()
}

func (s *subscribe) UnSubscribeReceiverKeys(receiverKey string, messageKeys ...string) {
	for _, key := range messageKeys {
		s.UnSubscribeReceiver(receiverKey, key)
	}
}

func (s *subscribe) FullPathString() string {
	return s.fullPath
}

func (s *subscribe) PathString() string {
	parts := strings.Split(s.fullPath, "/")
	return parts[len(parts)-1]
}

func (s *subscribe) SourceString() string {
	return strings.SplitN(s.fullPath, "/", 2)[0]
}
