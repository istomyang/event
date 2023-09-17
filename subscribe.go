package event

import (
	"fmt"
	"strings"
	"sync"
)

type subscribe struct {
	config                          *subscribeConfig
	fullPath                        string
	messageKeyHandlerMap            map[string]SubscribeHandler
	receiverKeyMessageKeyHandlerMap map[string]map[string]SubscribeHandler
	children                        []*subscribe // tree.
	mut                             sync.RWMutex
}

type subscribeConfig struct {
}

var _ Subscribe = &subscribe{}

func newSubscribe(config subscribeConfig) subscribe {
	return subscribe{
		config:                          &config,
		messageKeyHandlerMap:            make(map[string]SubscribeHandler),
		receiverKeyMessageKeyHandlerMap: make(map[string]map[string]SubscribeHandler),
		children:                        make([]*subscribe, 0),
	}
}

func (s *subscribe) Handle(context *Context) {
	// handle in this node layer.
	if context.Metadata.FullPath == s.FullPathString() {
		receiverKey := context.Metadata.ChannelKey
		messageKey := context.Metadata.Key
		var handler SubscribeHandler
		s.mut.RLocker()
		if receiverKey == "" {
			handler = s.messageKeyHandlerMap[messageKey]
		} else {
			handler = s.receiverKeyMessageKeyHandlerMap[receiverKey][messageKey]
		}
		s.mut.RUnlock()
		handler(context)
		return
	}
	// handle in next node layer.
	for _, child := range s.children {
		child.Handle(context)
	}
}

func (s *subscribe) Group(path string) Subscribe {
	var pre = s.fullPath
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	sub := &subscribe{
		config:                          s.config,
		fullPath:                        pre + "/" + path,
		messageKeyHandlerMap:            make(map[string]SubscribeHandler),
		receiverKeyMessageKeyHandlerMap: make(map[string]map[string]SubscribeHandler),
	}
	s.children = append(s.children, sub)
	return sub
}

func (s *subscribe) Subscribe(messageKey string, handler SubscribeHandler) {
	s.validateMessageKey(messageKey, s.messageKeyHandlerMap)
	s.mut.Lock()
	s.messageKeyHandlerMap[messageKey] = handler
	s.mut.Unlock()
}

func (s *subscribe) UnSubscribe(messageKey string) {
	s.mut.Lock()
	delete(s.messageKeyHandlerMap, messageKey)
	s.mut.Unlock()
}

func (s *subscribe) SubscribeReceiver(receiverKey string, messageKey string, handler SubscribeHandler) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if _, has := s.receiverKeyMessageKeyHandlerMap[receiverKey]; !has {
		s.receiverKeyMessageKeyHandlerMap[receiverKey] = make(map[string]SubscribeHandler)
	}
	m := s.receiverKeyMessageKeyHandlerMap[receiverKey]
	s.validateMessageKey(messageKey, m)
	m[messageKey] = handler
}

func (s *subscribe) UnSubscribeReceiver(receiverKey string, messageKey string) {
	s.mut.Lock()
	m := s.receiverKeyMessageKeyHandlerMap[receiverKey]
	delete(m, messageKey)
	s.mut.Unlock()
}

func (s *subscribe) validateMessageKey(messageKey string, m map[string]SubscribeHandler) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	if _, has := m[messageKey]; has {
		panic(fmt.Sprintf("ValidateMessageKey: duplicatie keys, got: %s", messageKey))
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
