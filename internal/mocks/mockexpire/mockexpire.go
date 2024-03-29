// Code generated by MockGen. DO NOT EDIT.
// Source: expire.go
//
// Generated by this command:
//
//	mockgen -source=expire.go -destination=../mocks/mockexpire/mockexpire.go -package=mockexpire
//

// Package mockexpire is a generated GoMock package.
package mockexpire

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockCacher is a mock of Cacher interface.
type MockCacher[K comparable, V any] struct {
	ctrl     *gomock.Controller
	recorder *MockCacherMockRecorder[K, V]
}

// MockCacherMockRecorder is the mock recorder for MockCacher.
type MockCacherMockRecorder[K comparable, V any] struct {
	mock *MockCacher[K, V]
}

// NewMockCacher creates a new mock instance.
func NewMockCacher[K comparable, V any](ctrl *gomock.Controller) *MockCacher[K, V] {
	mock := &MockCacher[K, V]{ctrl: ctrl}
	mock.recorder = &MockCacherMockRecorder[K, V]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacher[K, V]) EXPECT() *MockCacherMockRecorder[K, V] {
	return m.recorder
}

// Delete mocks base method.
func (m *MockCacher[K, V]) Delete(keys ...K) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Delete", varargs...)
}

// Delete indicates an expected call of Delete.
func (mr *MockCacherMockRecorder[K, V]) Delete(keys ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCacher[K, V])(nil).Delete), keys...)
}

// Keys mocks base method.
func (m *MockCacher[K, V]) Keys() []K {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Keys")
	ret0, _ := ret[0].([]K)
	return ret0
}

// Keys indicates an expected call of Keys.
func (mr *MockCacherMockRecorder[K, V]) Keys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockCacher[K, V])(nil).Keys))
}

// RandomKey mocks base method.
func (m *MockCacher[K, V]) RandomKey() (K, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomKey")
	ret0, _ := ret[0].(K)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// RandomKey indicates an expected call of RandomKey.
func (mr *MockCacherMockRecorder[K, V]) RandomKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomKey", reflect.TypeOf((*MockCacher[K, V])(nil).RandomKey))
}

// Size mocks base method.
func (m *MockCacher[K, V]) Size() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *MockCacherMockRecorder[K, V]) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockCacher[K, V])(nil).Size))
}

// TTL mocks base method.
func (m *MockCacher[K, V]) TTL(key K) (*time.Duration, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TTL", key)
	ret0, _ := ret[0].(*time.Duration)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// TTL indicates an expected call of TTL.
func (mr *MockCacherMockRecorder[K, V]) TTL(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TTL", reflect.TypeOf((*MockCacher[K, V])(nil).TTL), key)
}
