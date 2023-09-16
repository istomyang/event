// Package event
/*
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
*/
package event
