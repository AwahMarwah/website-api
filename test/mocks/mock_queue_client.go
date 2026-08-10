package mocks

import (
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/mock"
)

type MockQueueClient struct {
	mock.Mock
}

func (m *MockQueueClient) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	args := m.Called(task)

	var info *asynq.TaskInfo

	if args.Get(0) != nil {
		info = args.Get(0).(*asynq.TaskInfo)
	}

	return info, args.Error(1)
}
