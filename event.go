package event

type Publisher interface {
	Publish
	LifeCycle
}

type PublisherConfig struct {
	// Senders provides different options to send, if nil, use Default channel-based memory Sender.
	Senders map[string]Sender // senderKey:Sender
}

type Subscriber interface {
	Subscribe
	LifeCycle
}

type SubscriberConfig struct {
	// Receivers provides different options to receive, if nil, use Default channel-based memory Receiver.
	Receivers map[string]Receiver // receiverKey:Receiver
}

type Publish interface {
	Group(path string) Publish
	// Publish publishes all different Receiver.
	Publish(messageKey string, message Message)
	PublishSync(messageKey string, message Message) error
	PublishSender(senderKey string, messageKey string, message Message)
	PublishSyncSender(senderKey string, messageKey string, message Message) error
	Helper
}

type Subscribe interface {
	Group(path string) Subscribe
	// Subscribe will panic when duplicate messageKeys.
	Subscribe(messageKey string, handler SubscribeHandler) // Support wildcard.
	// SubscribeReceiver will panic when duplicate messageKeys.
	SubscribeReceiver(receiverKey string, messageKey string, handler SubscribeHandler)
	Helper
}

// Sender is a component doing really sending Message to Receiver.
// You can implement it by yourself or use default provided by this package.
type Sender interface {
	Send(Message)
	SendSync(Message) error
	LifeCycle
}

// Receiver is a component doing really receiving Message from Sender.
// You can implement it by yourself or use default provided by this package.
type Receiver interface {
	Receive() <-chan Message
	LifeCycle
}

type Context struct {
	*Message
}

type MessageMetadata struct {
	Key      string `json:"key" yaml:"key" bson:"key"`
	FullPath string `json:"full_path" yaml:"full_path" bson:"full_path"`
	Path     string `json:"path" yaml:"path" bson:"path"`
	// Source is first Path of FullPath.
	Source string `json:"source" yaml:"source" bson:"source"`

	// keyOfSenderAndReceiver is for specific channel should be same across Sender to Receiver.
	// Maybe empty string if call Publish.Publish.
	keyOfSenderAndReceiver string
}

type Message struct {
	// Metadata is Added by Publish.
	Metadata MessageMetadata `json:"metadata" yaml:"metadata" bson:"metadata"`
	Body     []byte          `json:"body" yaml:"body" bson:"body"`
	Extra    string          `json:"extra,omitempty" yaml:"extra,omitempty" bson:"extra,omitempty"`
}

type LifeCycle interface {
	Run() error
	MustRun()
	Close() error
}

type Helper interface {
	// FullPathString includes RootPath, MiddlePathPart and Path.
	FullPathString() string
	// PathString shows current Publish or Subscribe 's Path.
	PathString() string
	SourceString() string
}

type SubscribeHandler func(*Context)
