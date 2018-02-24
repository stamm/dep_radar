package mocks

import context "context"
import interfaces "github.com/stamm/dep_radar/interfaces"
import mock "github.com/stretchr/testify/mock"

// IProvider is an autogenerated mock type for the IProvider type
type IProvider struct {
	mock.Mock
}

// File provides a mock function with given fields: ctx, pkg, branch, filename
func (_m *IProvider) File(ctx context.Context, pkg interfaces.Pkg, branch string, filename string) ([]byte, error) {
	ret := _m.Called(ctx, pkg, branch, filename)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, interfaces.Pkg, string, string) []byte); ok {
		r0 = rf(ctx, pkg, branch, filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interfaces.Pkg, string, string) error); ok {
		r1 = rf(ctx, pkg, branch, filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GoGetURL provides a mock function with given fields:
func (_m *IProvider) GoGetURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Tags provides a mock function with given fields: _a0, _a1
func (_m *IProvider) Tags(_a0 context.Context, _a1 interfaces.Pkg) ([]interfaces.Tag, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []interfaces.Tag
	if rf, ok := ret.Get(0).(func(context.Context, interfaces.Pkg) []interfaces.Tag); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interfaces.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interfaces.Pkg) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

var _ interfaces.IProvider = (*IProvider)(nil)
