package pubsub

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type Subscriber struct {
	conn net.Conn
	ch   chan string
}
type PubSub struct {
	subscribers map[string][]*Subscriber
	mu          sync.Mutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]*Subscriber),
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

	ps.mu.Lock()
	subs := append([]*Subscriber{}, ps.subscribers[topic]...)
	ps.mu.Unlock()

	subs = ps.subscribers[topic]

	if subs == nil {
		return 404, errors.New("topic not found")
	}

	fmt.Println(subs)

	for _, i := range subs {
		fmt.Printf("Publishing message to topic '%s': %s\n", topic, msg)
		select {
		case i.ch <- msg:
		default:
			fmt.Println("message dropped for topic:", topic)
		}
	}
	return 200, nil
}
