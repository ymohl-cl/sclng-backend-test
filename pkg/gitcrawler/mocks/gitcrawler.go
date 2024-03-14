// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	gitcrawler "github.com/Scalingo/sclng-backend-test-v1/pkg/gitcrawler"
	mock "github.com/stretchr/testify/mock"
)

// GitCrawler is an autogenerated mock type for the GitCrawler type
type GitCrawler struct {
	mock.Mock
}

// EnrichRepositories provides a mock function with given fields: ctx, repositories
func (_m *GitCrawler) EnrichRepositories(ctx context.Context, repositories []gitcrawler.Repository) ([]gitcrawler.Repository, error) {
	ret := _m.Called(ctx, repositories)

	var r0 []gitcrawler.Repository
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []gitcrawler.Repository) ([]gitcrawler.Repository, error)); ok {
		return rf(ctx, repositories)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []gitcrawler.Repository) []gitcrawler.Repository); ok {
		r0 = rf(ctx, repositories)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]gitcrawler.Repository)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []gitcrawler.Repository) error); ok {
		r1 = rf(ctx, repositories)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublicRepository provides a mock function with given fields: ctx, nb, oldest
func (_m *GitCrawler) PublicRepository(ctx context.Context, nb int32, oldest bool) ([]gitcrawler.Repository, error) {
	ret := _m.Called(ctx, nb, oldest)

	var r0 []gitcrawler.Repository
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, bool) ([]gitcrawler.Repository, error)); ok {
		return rf(ctx, nb, oldest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, bool) []gitcrawler.Repository); ok {
		r0 = rf(ctx, nb, oldest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]gitcrawler.Repository)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, bool) error); ok {
		r1 = rf(ctx, nb, oldest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGitCrawler creates a new instance of GitCrawler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGitCrawler(t interface {
	mock.TestingT
	Cleanup(func())
}) *GitCrawler {
	mock := &GitCrawler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
