// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/goto/meteor/plugins/extractors/caramlstore/internal/core"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// CaraMLClient is an autogenerated mock type for the Client type
type CaraMLClient struct {
	mock.Mock
}

type CaraMLClient_Expecter struct {
	mock *mock.Mock
}

func (_m *CaraMLClient) EXPECT() *CaraMLClient_Expecter {
	return &CaraMLClient_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *CaraMLClient) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CaraMLClient_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type CaraMLClient_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *CaraMLClient_Expecter) Close() *CaraMLClient_Close_Call {
	return &CaraMLClient_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *CaraMLClient_Close_Call) Run(run func()) *CaraMLClient_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CaraMLClient_Close_Call) Return(_a0 error) *CaraMLClient_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

// Connect provides a mock function with given fields: ctx, host, maxSizeInMB, timeout
func (_m *CaraMLClient) Connect(ctx context.Context, host string, maxSizeInMB int, timeout time.Duration) error {
	ret := _m.Called(ctx, host, maxSizeInMB, timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, time.Duration) error); ok {
		r0 = rf(ctx, host, maxSizeInMB, timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CaraMLClient_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type CaraMLClient_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
//   - ctx context.Context
//   - host string
//   - maxSizeInMB int
//   - timeout time.Duration
func (_e *CaraMLClient_Expecter) Connect(ctx interface{}, host interface{}, maxSizeInMB interface{}, timeout interface{}) *CaraMLClient_Connect_Call {
	return &CaraMLClient_Connect_Call{Call: _e.mock.On("Connect", ctx, host, maxSizeInMB, timeout)}
}

func (_c *CaraMLClient_Connect_Call) Run(run func(ctx context.Context, host string, maxSizeInMB int, timeout time.Duration)) *CaraMLClient_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int), args[3].(time.Duration))
	})
	return _c
}

func (_c *CaraMLClient_Connect_Call) Return(_a0 error) *CaraMLClient_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

// Entities provides a mock function with given fields: ctx, project
func (_m *CaraMLClient) Entities(ctx context.Context, project string) (map[string]*core.Entity, error) {
	ret := _m.Called(ctx, project)

	var r0 map[string]*core.Entity
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]*core.Entity); ok {
		r0 = rf(ctx, project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*core.Entity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CaraMLClient_Entities_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Entities'
type CaraMLClient_Entities_Call struct {
	*mock.Call
}

// Entities is a helper method to define mock.On call
//   - ctx context.Context
//   - project string
func (_e *CaraMLClient_Expecter) Entities(ctx interface{}, project interface{}) *CaraMLClient_Entities_Call {
	return &CaraMLClient_Entities_Call{Call: _e.mock.On("Entities", ctx, project)}
}

func (_c *CaraMLClient_Entities_Call) Run(run func(ctx context.Context, project string)) *CaraMLClient_Entities_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CaraMLClient_Entities_Call) Return(_a0 map[string]*core.Entity, _a1 error) *CaraMLClient_Entities_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// FeatureTables provides a mock function with given fields: ctx, project
func (_m *CaraMLClient) FeatureTables(ctx context.Context, project string) ([]*core.FeatureTable, error) {
	ret := _m.Called(ctx, project)

	var r0 []*core.FeatureTable
	if rf, ok := ret.Get(0).(func(context.Context, string) []*core.FeatureTable); ok {
		r0 = rf(ctx, project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.FeatureTable)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CaraMLClient_FeatureTables_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FeatureTables'
type CaraMLClient_FeatureTables_Call struct {
	*mock.Call
}

// FeatureTables is a helper method to define mock.On call
//   - ctx context.Context
//   - project string
func (_e *CaraMLClient_Expecter) FeatureTables(ctx interface{}, project interface{}) *CaraMLClient_FeatureTables_Call {
	return &CaraMLClient_FeatureTables_Call{Call: _e.mock.On("FeatureTables", ctx, project)}
}

func (_c *CaraMLClient_FeatureTables_Call) Run(run func(ctx context.Context, project string)) *CaraMLClient_FeatureTables_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CaraMLClient_FeatureTables_Call) Return(_a0 []*core.FeatureTable, _a1 error) *CaraMLClient_FeatureTables_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Projects provides a mock function with given fields: ctx
func (_m *CaraMLClient) Projects(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CaraMLClient_Projects_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Projects'
type CaraMLClient_Projects_Call struct {
	*mock.Call
}

// Projects is a helper method to define mock.On call
//   - ctx context.Context
func (_e *CaraMLClient_Expecter) Projects(ctx interface{}) *CaraMLClient_Projects_Call {
	return &CaraMLClient_Projects_Call{Call: _e.mock.On("Projects", ctx)}
}

func (_c *CaraMLClient_Projects_Call) Run(run func(ctx context.Context)) *CaraMLClient_Projects_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *CaraMLClient_Projects_Call) Return(_a0 []string, _a1 error) *CaraMLClient_Projects_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewCaraMLClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewCaraMLClient creates a new instance of CaraMLClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCaraMLClient(t mockConstructorTestingTNewCaraMLClient) *CaraMLClient {
	mock := &CaraMLClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
