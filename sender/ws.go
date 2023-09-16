package sender

import (
	"context"
	"git.grizzlychina.com/infrastructures/event"
	"git.grizzlychina.com/infrastructures/event/pkg/ws"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type wsSender struct {
	config   *WebsocketConfig
	server   ws.Server
	sessions []ws.Session
	ctx      context.Context
	cancel   context.CancelFunc
	mut      sync.RWMutex
}

var _ event.Sender = &wsSender{}

type WebsocketConfig struct {
	Addr            string
	Path            string
	ReadBufferSize  int
	WriteBufferSize int
	Encode          func(event.Message) []byte
	CheckOrigin     func(r *http.Request) bool
}

func NewWebsocketSender(config WebsocketConfig) event.Sender {
	ctx, cancel := context.WithCancel(context.Background())
	return &wsSender{
		ctx:      ctx,
		cancel:   cancel,
		config:   &config,
		sessions: make([]ws.Session, 0),
		server: ws.NewServer(context.Background(), ws.ServerConfig{Upgrader: websocket.Upgrader{
			ReadBufferSize:  config.ReadBufferSize,
			WriteBufferSize: config.WriteBufferSize,
			CheckOrigin:     config.CheckOrigin,
		}}),
	}
}

func (w *wsSender) Send(message event.Message) {
	w.mut.RLocker()
	if w.sessions == nil {
		return
	}
	var localSessions = w.sessions
	w.mut.RUnlock()
	for _, session := range localSessions {
		session := session
		go func() {
			_ = session.Send(w.config.Encode(message))
		}()
	}
}

func (w *wsSender) SendSync(message event.Message) error {
	w.mut.RLocker()
	if w.sessions == nil {
		return nil
	}
	var localSessions = w.sessions
	w.mut.RUnlock()
	for _, session := range localSessions {
		if err := session.Send(w.config.Encode(message)); err != nil {
			return err
		}
	}
	return nil
}

func (w *wsSender) Run() error {
	http.HandleFunc(w.config.Path, func(writer http.ResponseWriter, reader *http.Request) {
		session, err := w.server.Create(writer, reader)
		if err != nil {
			return
		}
		w.mut.Lock()
		if w.sessions == nil {
			return
		}
		w.sessions = append(w.sessions, session)
		w.mut.Unlock()
		select {
		case <-w.ctx.Done():
			return
		}
	})
	go func() {
		select {
		case <-w.ctx.Done():
			return
		default:
			_ = http.ListenAndServe(w.config.Addr, nil)
		}
	}()
	return nil
}

func (w *wsSender) MustRun() {
	if err := w.Run(); err != nil {
		panic(err)
	}
}

func (w *wsSender) Close() error {
	w.mut.Lock()
	w.sessions = nil
	w.mut.Unlock()
	w.cancel()
	w.server.Close()
	return nil
}
