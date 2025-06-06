package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]any
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.constructors[name] = constructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructorFn, ok := c.constructors[name]
	if !ok {
		return nil, errors.New("constructor is not founded")
	}
	constructor, ok := constructorFn.(func() any)
	if ok {
		return constructor(), nil
	}
	return constructor, nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
