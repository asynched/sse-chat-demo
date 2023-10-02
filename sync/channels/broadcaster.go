package channels

type Broadcaster[T any] struct {
	channels []chan T
}

func (broadcaster *Broadcaster[T]) Subscribe() chan T {
	channel := make(chan T)
	broadcaster.channels = append(broadcaster.channels, channel)

	return channel
}

func (broadcaster *Broadcaster[T]) Unsubscribe(channel chan T) {
	close(channel)

	for i, c := range broadcaster.channels {
		if c == channel {
			broadcaster.channels = append(broadcaster.channels[:i], broadcaster.channels[i+1:]...)
			break
		}
	}
}

func (broadcaster *Broadcaster[T]) Broadcast(message T) {
	for _, channel := range broadcaster.channels {
		channel <- message
	}
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		channels: make([]chan T, 0),
	}
}
