package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type NumType interface {
	int | int8 | int16 | int32 | int64
}

type CircularQueue[T NumType] struct {
	values []T
	// Индекс будущего элемента в слайсе
	nextValIdx int
	// Индекс следующего к удалению элемента
	nextToDelIdx int
	// Размер массива без нулей
	size int
}

func NewCircularQueue[T NumType](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
	} // need to implement
}

func (q *CircularQueue[T]) Push(value T) bool {
	if !q.Full() {
		q.values[q.nextValIdx] = value
		q.nextValIdx = (q.nextValIdx + 1) % len(q.values)
		q.size++
		return true
	}
	return false // need to implement
}

func (q *CircularQueue[T]) Pop() bool {
	if !q.Empty() {
		q.values[q.nextToDelIdx] = 0
		q.nextToDelIdx = (q.nextToDelIdx + 1) % len(q.values)
		q.size--
		return true
	}
	return false // need to implement
}

func (q *CircularQueue[T]) Front() T {
	if !q.Empty() {
		return q.values[q.nextToDelIdx]
	}
	return -1 // need to implement
}

func (q *CircularQueue[T]) Back() T {
	if !q.Empty() {
		if q.nextValIdx == 0 {
			return q.values[len(q.values)-1]
		}
		return q.values[q.nextValIdx-1]
	}
	return -1 // need to implement
}

func (q *CircularQueue[T]) Empty() bool {
	if q.size == 0 {
		return true
	}
	return false // need to implement
}

func (q *CircularQueue[T]) Full() bool {
	if q.size == len(q.values) {
		return true
	}
	return false // need to implement
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

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

	//	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
