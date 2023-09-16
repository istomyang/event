package event

type publisher struct {
	config *PublisherConfig
	publish
}

func NewPublisher(config PublisherConfig) Publisher {
	var publishConfig publishConfig
	if config.Senders == nil {
		if defaultMemoryCh != nil {
			defaultMemoryCh = make(chan Message)
		}
		publishConfig.Senders = map[string]Sender{defaultSenderKey: createDefaultSender()}
	}
	return &publisher{
		config:  &config,
		publish: newPublish(publishConfig),
	}
}

var _ Publisher = &publisher{}

func (p *publisher) Run() error {
	var err error
	var senders = p.config.Senders
	for _, sender := range senders {
		if err = sender.Run(); err != nil {
			return err
		}
	}
	return err
}

func (p *publisher) MustRun() {
	var senders = p.config.Senders
	for _, sender := range senders {
		sender.MustRun()
	}
}

func (p *publisher) Close() error {
	var err error
	var senders = p.config.Senders
	for _, sender := range senders {
		if err = sender.Close(); err != nil {
			return err
		}
	}
	return err
}
