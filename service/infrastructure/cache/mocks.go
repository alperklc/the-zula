package cache

import (
	"github.com/stretchr/testify/mock"
)

type MockedCache[T any] struct {
	mock.Mock
}

func (m *MockedCache[T]) Reset(id string) {
	m.Called(id)
}

func (m *MockedCache[T]) Write(id string, obj T) {
	m.Called(id, obj)
}

func (m *MockedCache[T]) Read(id string) *T {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*T)
}

var _ CacheInterface[any] = (*MockedCache[any])(nil)
