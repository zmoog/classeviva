// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
	spaggiari "github.com/zmoog/classeviva/adapters/spaggiari"
)

// Fetcher is an autogenerated mock type for the Fetcher type
type Fetcher struct {
	mock.Mock
}

// Fetch provides a mock function with given fields:
func (_m *Fetcher) Fetch() (spaggiari.Identity, error) {
	ret := _m.Called()

	var r0 spaggiari.Identity
	if rf, ok := ret.Get(0).(func() spaggiari.Identity); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(spaggiari.Identity)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFetcher creates a new instance of Fetcher. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFetcher(t testing.TB) *Fetcher {
	mock := &Fetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}