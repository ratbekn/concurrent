package queue

import (
	"fmt"
	"sync"
)

type Queue[T any] struct {
	head       int
	tail       int
	storage    []T
	lock       sync.Mutex
	isNotEmpty *sync.Cond
}

func New[T any](capacity int) *Queue[T] {
	q := &Queue[T]{head: -1, tail: -1, storage: make([]T, capacity)}
	q.isNotEmpty = sync.NewCond(&q.lock)

	return q
}

func (q *Queue[T]) isEmpty() bool {
	return q.head == -1 && q.tail == -1
}

func (q *Queue[T]) isFull() bool {
	return (q.tail+1)%len(q.storage) == q.head
}

func (q *Queue[T]) Push(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.isEmpty() {
		q.head = 0
	}
	q.tail = (q.tail + 1) % len(q.storage)
	q.storage[q.tail] = item

	q.isNotEmpty.Signal()
}

func (q *Queue[T]) Pop() T {
	q.lock.Lock()
	defer q.lock.Unlock()

	for q.isEmpty() {
		q.isNotEmpty.Wait()
	}

	item := q.storage[q.head]
	if q.head == q.tail {
		q.head = -1
		q.tail = -1

		return item
	}

	q.head = (q.head + 1) % len(q.storage)

	return item
}

func (q *Queue[T]) Print() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head <= q.tail {
		fmt.Println(q.storage[q.head : q.tail+1])
	} else {
		fmt.Println(append(q.storage[q.head:], q.storage[:q.tail+1]...))
	}
}
