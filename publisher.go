package pubsub

// An Event should have a Type, that's an integer (iota)
// and optionally an additional data object
type Event struct {
	Type int
	Data interface{}
}

// Listeners for the various types of events
var listeners map[int][]chan *Event = make(map[int][]chan *Event)

// Notify all listeners of this type of event
func Publish(eventType int, data interface{}) {
	if listeners[eventType] != nil {
		event := &Event{eventType, data}
		for _, c := range listeners[eventType] {
			c <- event
		}
	}
}

// Subscribe to an specific type of Event and pass a channel you might
// want to reuse to receive the notifications, or nil to get a new one
// from the library
func Subscribe(eventType int, channel chan *Event) chan *Event {
	if channel == nil {
		channel = make(chan *Event)
	}
	listeners[eventType] = append(listeners[eventType], channel)
	return channel
}

// Unsubscribe from the event
func Unsubscribe(eventType int, channel chan *Event) bool {
	if listeners[eventType] != nil {
		for i, c := range listeners[eventType] {
			if channel == c {
				listeners[eventType] = append(listeners[eventType][:i], listeners[eventType][i+1:]...)
				return true
			}
		}
	}
	return false
}
