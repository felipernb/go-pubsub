go-pubsub
=========

A very basic pubsub mechanism using goroutines and channels

Author: Felipe Ribeiro <felipernb@gmail.com>

Usage
-----

```
$ go get github.com/felipernb/go-pubsub
```

myfile.go
```go
package main

import (
	"github.com/felipernb/go-pubsub"
	"fmt"
)

const (
	EVENT_TYPE0 = iota
	EVENT_TYPE1 = iota
)

func main() {
	// Subscribe to an event with a brand new channel
	c := pubsub.Subscribe(EVENT_TYPE0, nil)
	c1 := pubsub.Subscribe(EVENT_TYPE1, nil)


	// Subscribe to an event reusing a channel
	pubsub.Subscribe(EVENT_TYPE1, c)

	// Listen for events
	go func() {
		var event *pubsub.Event
		for {
			select {
			case event = <-c:
				fmt.Printf("Event type %d on channel C with data: %v\n", event.Type, event.Data)
			case event = <-c1:
				fmt.Printf("Event type %d on channel C1 with data: %v\n", event.Type, event.Data)
			}
		}
	}()

	// Publish events
	pubsub.Publish(EVENT_TYPE1, "Event 1 data")
	pubsub.Publish(EVENT_TYPE0, "Event 0 data")
}

```
