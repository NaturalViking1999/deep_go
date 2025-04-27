package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node[T cmp.Ordered] struct {
	left, right *Node[T]
	key         T
	value       any
}

type OrderedMap[T cmp.Ordered] struct {
	node *Node[T]
	size int
}

func NewOrderedMap[T cmp.Ordered]() OrderedMap[T] {
	return OrderedMap[T]{}
}

func (m *OrderedMap[T]) Insert(key T, value any) {
	x := m.node
	var y *Node[T] = nil

	for x != nil {
		compared := compare(key, x.key)
		if compared == 0 {
			x.value = value
			return
		} else {
			y = x
			if compared < 0 {
				x = x.left
			} else {
				x = x.right
			}
		}
	}

	newNode := &Node[T]{key: key, value: value}
	if y == nil {
		m.node = newNode
	} else {
		compared := compare(key, y.key)
		if compared < 0 {
			y.left = newNode
		} else {
			y.right = newNode
		}
	}
	m.size++
}

func compare[T cmp.Ordered](k1, k2 T) int {
	if k1 == k2 {
		return 0
	}
	if k1 < k2 {
		return -1
	}
	return 1
}

func (m *OrderedMap[T]) Erase(key T) {
	x := m.node
	y := &Node[T]{}

	for x != nil {
		compared := compare(key, x.key)
		if compared == 0 {
			break
		} else {
			y = x
			if compared < 0 {
				x = x.left
			} else {
				x = x.right
			}
		}
	}

	if x == nil {
		return
	}

	if x.right == nil {
		if y == nil {
			m.node = x.left
		} else {
			if x == y.left {
				y.left = x.left
			} else {
				y.right = x.left
			}
		}
	} else {
		leftMost := x.right
		y = nil
		for leftMost.left != nil {
			y = leftMost
			leftMost = leftMost.left
		}

		if y != nil {
			y.left = leftMost.right
		} else {
			x.right = leftMost.right
		}

		x.key = leftMost.key
		x.value = leftMost.value
	}
	m.size--
}

func (m *OrderedMap[T]) Contains(key T) bool {
	x := m.node
	var compared int
	for x != nil {
		compared = compare(key, x.key)
		if compared == 0 {
			return true
		} else if compared < 0 {
			x = x.left
		} else {
			x = x.right
		}
	}
	return false
}

func (m *OrderedMap[T]) Size() int {
	return m.size
}

func (m *OrderedMap[T]) ForEach(action func(T, any)) {
	x := m.node
	if x != nil {
		forEach(action, m.node)
	}
}

func forEach[T cmp.Ordered](action func(T, any), x *Node[T]) {
	if x.left != nil {
		forEach(action, x.left)
	}
	action(x.key, x.value)
	if x.right != nil {
		forEach(action, x.right)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
