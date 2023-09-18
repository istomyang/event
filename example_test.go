package event

import (
	"fmt"
	"time"
)

func Example() {
	runPublisher()
	runSubscriber()

	t := time.After(time.Second * 3)
	<-t

	// Output:
}

func runPublisher() {

	var publisher = NewPublisher(PublisherConfig{})

	// Attach lifecycle to App's lifecycle.
	publisher.MustRun()
	defer publisher.Close()

	// Root
	rootModule := publisher.Group("rootModule")
	{
		rootModule.Publish("path1", Message{
			Body:  []byte(fmt.Sprintf("Publisher: %s: Ping!", rootModule.FullPathString())),
			Extra: "",
		})

		// Sub
		subModule := rootModule.Group("subModule")
		{
			subModule.Publish("path2", Message{
				Body:  []byte(fmt.Sprintf("Publisher: %s: Ping!", subModule.FullPathString())),
				Extra: "",
			})

			// SubSub
			subSubModule := subModule.Group("subSubModule")
			{
				subSubModule.Publish("path3", Message{
					Body:  []byte(fmt.Sprintf("Publisher: %s: Ping!", subSubModule.FullPathString())),
					Extra: "",
				})
			}
		}
	}
}

func runSubscriber() {
	var subscriber = NewSubscriber(SubscriberConfig{})

	// Attach lifecycle to App's lifecycle.
	subscriber.MustRun()
	defer subscriber.Close()

	// Root
	rootModule := subscriber.Group("rootModule")
	{
		rootModule.Subscribe("path1", func(context *Context) {
			fmt.Printf("Receiver: %s [%s]\n", rootModule.FullPathString(), string(context.Message.Body))
		})

		// Sub
		subModule := rootModule.Group("subModule")
		{
			subModule.Subscribe("path2", func(context *Context) {
				fmt.Printf("Receiver: %s [%s]\n", subModule.FullPathString(), string(context.Message.Body))
			})

			// SubSub
			subSubModule := subModule.Group("subSubModule")
			{
				subSubModule.Subscribe("path3", func(context *Context) {
					fmt.Printf("Receiver: %s [%s]\n", subSubModule.FullPathString(), string(context.Message.Body))
				})
			}
		}
	}
}
