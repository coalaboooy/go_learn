package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	size   int
	index  int // index указывает на текущее положение (курсор) в очереди
}

func NewCircularQueue(size int) CircularQueue {
	cq := CircularQueue{
		make([]int, size),
		0,
		0,
	}
	return cq
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}

	q.values[q.index] = value
	q.index = (q.index + 1) % (len(q.values))
	q.size += 1
	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	q.values[0] = 0
	q.index = 0
	q.size -= 1
	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.index]
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	return q.values[(q.index+len(q.values)-1)%(len(q.values))]
}

func (q *CircularQueue) Empty() bool {
	return q.size == 0
}

func (q *CircularQueue) Full() bool {
	return q.size == len(q.values)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	// Тесты с другими длинами очереди

	// длина = 1
	oneQueue := NewCircularQueue(1)

	assert.True(t, oneQueue.Empty())
	assert.False(t, oneQueue.Full())

	assert.Equal(t, -1, oneQueue.Front())
	assert.Equal(t, -1, oneQueue.Back())
	assert.False(t, oneQueue.Pop())

	assert.True(t, oneQueue.Push(1))
	assert.False(t, oneQueue.Push(5))

	assert.True(t, reflect.DeepEqual([]int{1}, oneQueue.values))

	assert.False(t, oneQueue.Empty())
	assert.True(t, oneQueue.Full())

	assert.Equal(t, 1, oneQueue.Front())
	assert.Equal(t, 1, oneQueue.Back())

	assert.True(t, oneQueue.Pop())
	assert.True(t, oneQueue.Empty())
	assert.False(t, oneQueue.Full())
	assert.True(t, oneQueue.Push(4))

	assert.True(t, oneQueue.Pop())
	assert.False(t, oneQueue.Pop())

	assert.True(t, oneQueue.Empty())
	assert.False(t, oneQueue.Full())

	// длина = 4
	fourQueue := NewCircularQueue(4)

	assert.True(t, fourQueue.Empty())
	assert.False(t, fourQueue.Full())

	assert.Equal(t, -1, fourQueue.Front())
	assert.Equal(t, -1, fourQueue.Back())
	assert.False(t, fourQueue.Pop())

	assert.True(t, fourQueue.Push(5))
	assert.True(t, fourQueue.Push(10))
	assert.True(t, fourQueue.Push(15))
	assert.True(t, fourQueue.Push(20))
	assert.False(t, fourQueue.Push(9))

	assert.True(t, reflect.DeepEqual([]int{5, 10, 15, 20}, fourQueue.values))

	assert.False(t, fourQueue.Empty())
	assert.True(t, fourQueue.Full())

	assert.Equal(t, 5, fourQueue.Front())
	assert.Equal(t, 20, fourQueue.Back())

	assert.True(t, fourQueue.Pop())
	assert.True(t, fourQueue.Pop())
	assert.True(t, fourQueue.Pop())
	assert.True(t, fourQueue.Pop())
	assert.True(t, fourQueue.Empty())
	assert.False(t, fourQueue.Full())
	assert.True(t, fourQueue.Push(9))

	assert.True(t, fourQueue.Pop())
	assert.False(t, fourQueue.Pop())

	assert.True(t, fourQueue.Empty())
	assert.False(t, fourQueue.Full())

}
