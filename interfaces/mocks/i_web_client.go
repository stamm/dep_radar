package mocks

import interfaces "github.com/stamm/dep_radar/interfaces"
import mock "github.com/stretchr/testify/mock"

// IWebClient is an autogenerated mock type for the IWebClient type
type IWebClient struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *IWebClient) Get(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

var _ interfaces.IWebClient = (*IWebClient)(nil)