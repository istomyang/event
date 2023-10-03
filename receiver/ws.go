package receiver

import (
	"context"
	"github.com/istomyang/event"
	"github.com/istomyang/event/pkg/ws"
	"sync/atomic"
)

type wsReceiver struct {
	config      *WebsocketConfig
	session     ws.Session
	messageChan chan event.Message
	client      ws.Client
	chanClosed  atomic.Bool
}

var _ event.Receiver = &wsReceiver{}

type WebsocketConfig struct {
	Addr   string
	Path   string
	Decode func([]byte) event.Message
}

func NewWebsocketReceiver(config WebsocketConfig) event.Receiver {
	return &wsReceiver{
		config: &config,
		client: ws.NewClient(context.Background(), ws.ClientConfig{}),
	}
}

func (w *wsReceiver) Receive() <-chan event.Message {
	return w.messageChan
}

func (w *wsReceiver) Run() error {
	var err error
	w.session, err = w.client.Create(w.config.Addr, w.config.Path)
	go func() {
		for data := range w.session.Receive() {
			if w.chanClosed.Load() {
				break
			}
			w.messageChan <- w.config.Decode(data)
		}
	}()
	return err
}

func (w *wsReceiver) MustRun() {
	if err := w.Run(); err != nil {
		panic(err)
	}
}

func (w *wsReceiver) Close() error {
	var err error
	w.session = nil
	w.client.Close()
	w.chanClosed.Store(true)
	close(w.messageChan)
	return err
}
