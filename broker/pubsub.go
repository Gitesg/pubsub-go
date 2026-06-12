package pubsub

import (
	"fmt"
	"net"
	"sync"
)

type Subscriber struct {
	conn net.Conn
	ch   chan string
}

type Message struct {
	topic   string
	message string
}
type PubSub struct {
	subscribers map[string][]*Subscriber
	mu          sync.Mutex
	globalqueue chan Message
}

func NewPubSub() *PubSub {

	ps := &PubSub{
		subscribers: make(map[string][]*Subscriber),
		mu:          sync.Mutex{},
		globalqueue: make(chan Message, 100),
	}
	ps.Dispatcher(3)
	return ps
}

func (ps *PubSub) Dispatcher(worker int) {
	for i := 0; i < worker; i++ {
		go func() {
			for msg := range ps.globalqueue {

				topic, mx := msg.topic, msg.message

				ps.mu.Lock()
				subs, ok := ps.subscribers[topic]
				if !ok {
					ps.mu.Unlock()
					continue
				}
				cop1 := make([]*Subscriber, len(subs))
				copy(cop1, subs)
				ps.mu.Unlock()

				for _, sub := range cop1 {
					select {
					case sub.ch <- mx:
					default:
						fmt.Println("slow subscriber, dropping")
					}
				}
			}
		}()
	}

}

func (ps *PubSub) Subscribe(topic string, conn net.Conn) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 10)
	ps.subscribers[topic] = append(ps.subscribers[topic], &Subscriber{
		conn: conn,
		ch:   ch,
	})
	return ch
}

func (ps *PubSub) Publisher(msg string, topic string) (code int, e error) {

	select {
	case ps.globalqueue <- Message{topic: topic, message: msg}:
	default:
		return 0, fmt.Errorf("mqeue is full")
	}

	return 200, nil

}
