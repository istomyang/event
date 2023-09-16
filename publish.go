package event

import "strings"

type publish struct {
	config   *publishConfig
	fullPath string
}

type publishConfig struct {
	Senders map[string]Sender
}

func newPublish(config publishConfig) publish {
	return publish{config: &config}
}

var _ Publish = &publish{}

func (p *publish) Group(path string) Publish {
	var pre = p.fullPath
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	return &publish{
		config:   p.config,
		fullPath: pre + "/" + path,
	}
}

func (p *publish) Publish(messageKey string, message Message) {
	message.Metadata = p.createMessageMetadata(messageKey)
	for _, s := range p.config.Senders {
		s := s
		go func() {
			s.Send(message)
		}()
	}
}

func (p *publish) PublishSync(messageKey string, message Message) error {
	message.Metadata = p.createMessageMetadata(messageKey)
	for _, s := range p.config.Senders {
		if err := s.SendSync(message); err != nil {
			return err
		}
	}
	return nil
}

func (p *publish) PublishSender(senderKey string, messageKey string, message Message) {
	message.Metadata = p.createMessageMetadata(messageKey)
	message.Metadata.ChannelKey = senderKey
	s := p.config.Senders[senderKey]
	s.Send(message) // s is ok when hasn't.
}

func (p *publish) PublishSyncSender(senderKey string, messageKey string, message Message) error {
	message.Metadata = p.createMessageMetadata(messageKey)
	message.Metadata.ChannelKey = senderKey
	s := p.config.Senders[senderKey]
	return s.SendSync(message) // s is ok when hasn't.
}

func (p *publish) FullPathString() string {
	return p.fullPath
}

func (p *publish) PathString() string {
	parts := strings.Split(p.fullPath, "/")
	return parts[len(parts)-1]
}

func (p *publish) SourceString() string {
	return strings.SplitN(p.fullPath, "/", 2)[0]
}

func (p *publish) createMessageMetadata(messageKey string) MessageMetadata {
	return MessageMetadata{
		FullPath: p.FullPathString(),
		Path:     p.PathString(),
		Source:   p.SourceString(),
		Key:      messageKey,
	}
}
