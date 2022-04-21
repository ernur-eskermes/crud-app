package otp

import "github.com/stretchr/testify/mock"

type MockGenerator struct {
	mock.Mock
}

func (m *MockGenerator) RandomSecret(length int) string {
	args := m.Called(length)

	v, _ := args.Get(0).(string)

	return v
}
