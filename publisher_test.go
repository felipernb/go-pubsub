package pubsub

import (
	"fmt"
	"testing"
	"time"
)

const (
	EVENT_FOLLOW   = iota // An event that is triggered when a user follows another
	EVENT_UNFOLLOW = iota // Event that is triggered when a user unfollows another
)

func TestPubSub(t *testing.T) {
	followEvents := 0
	unfollowEvents := 0
	followNotify := 0

	follow := Subscribe(EVENT_FOLLOW, nil)
	follow2 := Subscribe(EVENT_FOLLOW, nil)
	unfollow := Subscribe(EVENT_UNFOLLOW, nil)

	go func() {
		var e *Event
		for {
			select {
			case e = <-follow:
				if e.Type != EVENT_FOLLOW {
					fmt.Println("Expected EVENT_FOLLOW, got %v instead", e.Type)
					t.Fail()
				}
				followEvents++
				followNotify++
			case e = <-unfollow:
				if e.Type != EVENT_UNFOLLOW {
					fmt.Println("Expected EVENT_UNFOLLOW, got %v instead", e.Type)
					t.Fail()
				}
				unfollowEvents++
			case e = <-follow2:
				followNotify++
			}
		}
	}()

	Publish(EVENT_FOLLOW, nil)
	Publish(EVENT_UNFOLLOW, nil)
	Publish(EVENT_FOLLOW, nil)

	// One second for all events to be handled
	time.Sleep(time.Second)

	if followEvents != 2 {
		fmt.Println("2 FOLLOW events should had been triggered")
		t.Fail()
	}

	if followNotify != 4 {
		fmt.Printf("2 FOLLOW events should had been notified to 2 handlers, it was notified %d times instead\n", followNotify)
		t.Fail()
	}

	if unfollowEvents != 1 {
		fmt.Println("1 UNFOLLOW event should had been triggered")
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	ch := make(chan *Event)

	if Unsubscribe(EVENT_FOLLOW, ch) {
		fmt.Println("The channel has not being subscribed yet, shouldn't be able to unsubscribe")
		t.Fail()
	}

	ch2 := Subscribe(EVENT_FOLLOW, nil)

	if Unsubscribe(EVENT_FOLLOW, ch) {
		fmt.Println("The channel has not being subscribed yet, shouldn't be able to unsubscribe")
		t.Fail()
	}

	Subscribe(EVENT_FOLLOW, ch)

	eventHandleCounter1 := 0
	eventHandleCounter2 := 0
	timeout := time.After(time.Second)
	go func() {
		for {
			select {
			case <-ch:
				eventHandleCounter1++
			case <-ch2:
				eventHandleCounter2++
			case <-timeout:
				t.Fail()
			}

			if eventHandleCounter1 == 1 && eventHandleCounter2 == 1 {
				break
			}
		}
	}()
	Publish(EVENT_FOLLOW, nil)

	if !Unsubscribe(EVENT_FOLLOW, ch) {
		fmt.Println("The channel has being subscribed, should be able to unsubscribe")
		t.Fail()
	}

	eventHandleCounter2 = 0
	timeout = time.After(time.Second)
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println("The channel has been unsubscribed, shouldn't receive another message")
				t.Fail()
			case <-ch2:
				eventHandleCounter2++
			case <-timeout:
				if eventHandleCounter2 != 1 {
					fmt.Println("Event should had been handled once")
					t.Fail()
				}
				return
			}

		}
	}()

	Publish(EVENT_FOLLOW, nil)

}
