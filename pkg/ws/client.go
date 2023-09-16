package ws

import (
	"context"
	"github.com/gorilla/websocket"
)

// Client manages WebSocket connections in client end.
type Client interface {
	// Create creates a connection over a http connection and return a Session object.
	Create(addr string, path string) (Session, error)

	Run()
	Close()
}

type client struct {
	ctx      context.Context
	cancel   context.CancelFunc
	config   ClientConfig
	sessions []innerSession
}

func NewClient(ctx context.Context, config ClientConfig) Client {
	ctx, cancel := context.WithCancel(ctx)
	return &client{
		ctx:      ctx,
		cancel:   cancel,
		config:   config,
		sessions: make([]innerSession, 0), // buffer has data loss when panicked.
	}
}

func (c *client) Create(addr string, path string) (Session, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(c.ctx, addr+path, nil)
	if err != nil {

		return nil, err
	}

	var se = newSession(c.ctx, conn, path)
	se.(innerSession).Attach()
	c.sessions = append(c.sessions, se.(innerSession))

	return se, nil
}

func (c *client) Run() {
	go func() {
		select {
		case <-c.ctx.Done():

			c.Close()
		}
	}()
}

func (c *client) Close() {
	defer c.cancel()
	for _, se := range c.sessions {
		se.Close()
	}

}

var _ Client = &client{}

type fakeClient struct {
	config FakeClientConfig
}

type FakeClientConfig struct {
	ClientSend <-chan []byte
}

func NewFakeClient(config FakeClientConfig) Client {
	return &fakeClient{config: config}
}

func (f *fakeClient) Create(addr string, path string) (Session, error) {

	return newFakeSession(FakeSessionConfig{ClientSend: f.config.ClientSend}), nil
}

func (f *fakeClient) Run() {

}

func (f *fakeClient) Close() {

}

var _ Client = &fakeClient{}
