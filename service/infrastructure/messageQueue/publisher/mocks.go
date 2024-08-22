package mqpublisher

import (
	"github.com/stretchr/testify/mock"
)

type MockedMessagePublisher struct {
	mock.Mock
}

func (m *MockedMessagePublisher) Publish(body ActivityMessage) error {
	args := m.Called(body)

	return args.Error(0)
}
