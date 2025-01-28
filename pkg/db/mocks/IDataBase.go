// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// IDataBase is an autogenerated mock type for the IDataBase type
type IDataBase struct {
	mock.Mock
}

// CreateCollection provides a mock function with given fields: name
func (_m *IDataBase) CreateCollection(name string) *mongo.Collection {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for CreateCollection")
	}

	var r0 *mongo.Collection
	if rf, ok := ret.Get(0).(func(string) *mongo.Collection); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.Collection)
		}
	}

	return r0
}

// GetCollection provides a mock function with given fields: name
func (_m *IDataBase) GetCollection(name string) *mongo.Collection {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetCollection")
	}

	var r0 *mongo.Collection
	if rf, ok := ret.Get(0).(func(string) *mongo.Collection); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.Collection)
		}
	}

	return r0
}

// GetDatabase provides a mock function with no fields
func (_m *IDataBase) GetDatabase() *mongo.Database {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDatabase")
	}

	var r0 *mongo.Database
	if rf, ok := ret.Get(0).(func() *mongo.Database); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.Database)
		}
	}

	return r0
}

// NewIDataBase creates a new instance of IDataBase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDataBase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDataBase {
	mock := &IDataBase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
