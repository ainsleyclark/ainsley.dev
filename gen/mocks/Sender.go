// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Sender is an autogenerated mock type for the Sender type
type Sender struct {
	mock.Mock
}

// Send provides a mock function with given fields: ctx, channelID, subject, message
func (_m *Sender) Send(ctx context.Context, channelID string, subject string, message string) error {
	ret := _m.Called(ctx, channelID, subject, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, channelID, subject, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSender interface {
	mock.TestingT
	Cleanup(func())
}

// NewSender creates a new instance of Sender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSender(t mockConstructorTestingTNewSender) *Sender {
	mock := &Sender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
