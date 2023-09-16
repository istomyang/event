# Event

Event is a library between the message transport implementation and application layer.

```
                                 ┌─────────────┐
                                 │  Publisher  │
                                 └──────┬──────┘
                                        │
                                        ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                                                                                  │
│  ┌────────────┐ ┌─────────────────────┐ ┌─────────────────┐ ┌─────────────────┐  │
│  │ Go Channel │ │ Kafka and Other MQs │ │ Websocket/Http2 │ │  RPC Framework  │  │
│  └────────────┘ └─────────────────────┘ └─────────────────┘ └─────────────────┘  │
│                                                                                  │
│                                     Broker                                       │
└───────────────────────────────────────┬──────────────────────────────────────────┘
                                        │
                                        ▼
                                 ┌──────────────┐
                                 │  Subscriber  │
                                 └──────────────┘
```