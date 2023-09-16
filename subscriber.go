package event

type subscriber struct {
	config    *SubscriberConfig
	subscribe // A root node of a tree.
}

func NewSubscriber(config SubscriberConfig) Subscriber {
	if config.Receivers == nil {
		config.Receivers = map[string]Receiver{defaultReceiverKey: createDefaultReceiver()}
	}
	return &subscriber{
		config:    &config,
		subscribe: newSubscribe(subscribeConfig{}),
	}
}

var _ Subscriber = &subscriber{}

func (s *subscriber) Run() error {
	var err error
	var receivers = s.config.Receivers
	for _, receiver := range receivers {
		if err = receiver.Run(); err != nil {
			return err
		}
	}
	s.launch()
	return err
}

func (s *subscriber) MustRun() {
	var receivers = s.config.Receivers
	for _, receiver := range receivers {
		receiver.MustRun()
	}
	s.launch()
}

func (s *subscriber) Close() error {
	var err error
	var receivers = s.config.Receivers
	for _, receiver := range receivers {
		if err = receiver.Close(); err != nil {
			return err
		}
	}
	return err
}

func (s *subscriber) launch() {
	for _, rec := range s.config.Receivers {
		rec := rec
		go func() {
			for message := range rec.Receive() {
				s.subscribe.Handle(&Context{&message})
			}
		}()
	}
}
