// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/wafer-bw/memcache/internal/expire (interfaces: IntCacher)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/expire/expire.go -package=mockexpire . IntCacher
//

// Package mockexpire is a generated GoMock package.
package mockexpire

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockIntCacher is a mock of IntCacher interface.
type MockIntCacher struct {
	ctrl     *gomock.Controller
	recorder *MockIntCacherMockRecorder
}

// MockIntCacherMockRecorder is the mock recorder for MockIntCacher.
type MockIntCacherMockRecorder struct {
	mock *MockIntCacher
}

// NewMockIntCacher creates a new mock instance.
func NewMockIntCacher(ctrl *gomock.Controller) *MockIntCacher {
	mock := &MockIntCacher{ctrl: ctrl}
	mock.recorder = &MockIntCacherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIntCacher) EXPECT() *MockIntCacherMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockIntCacher) Delete(arg0 ...int) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Delete", varargs...)
}

// Delete indicates an expected call of Delete.
func (mr *MockIntCacherMockRecorder) Delete(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIntCacher)(nil).Delete), arg0...)
}

// Keys mocks base method.
func (m *MockIntCacher) Keys() []int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Keys")
	ret0, _ := ret[0].([]int)
	return ret0
}

// Keys indicates an expected call of Keys.
func (mr *MockIntCacherMockRecorder) Keys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockIntCacher)(nil).Keys))
}

// RandomKey mocks base method.
func (m *MockIntCacher) RandomKey() (int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomKey")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// RandomKey indicates an expected call of RandomKey.
func (mr *MockIntCacherMockRecorder) RandomKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomKey", reflect.TypeOf((*MockIntCacher)(nil).RandomKey))
}

// Size mocks base method.
func (m *MockIntCacher) Size() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *MockIntCacherMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockIntCacher)(nil).Size))
}

// TTL mocks base method.
func (m *MockIntCacher) TTL(arg0 int) (*time.Duration, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TTL", arg0)
	ret0, _ := ret[0].(*time.Duration)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// TTL indicates an expected call of TTL.
func (mr *MockIntCacherMockRecorder) TTL(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TTL", reflect.TypeOf((*MockIntCacher)(nil).TTL), arg0)
}
