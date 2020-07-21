// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/msg/producer/writer/consumer_service_go

// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package writer is a generated GoMock package.
package writer

import (
	"reflect"

	"github.com/m3db/m3/src/msg/producer"

	"github.com/golang/mock/gomock"
)

// MockconsumerServiceWriter is a mock of consumerServiceWriter interface
type MockconsumerServiceWriter struct {
	ctrl     *gomock.Controller
	recorder *MockconsumerServiceWriterMockRecorder
}

// MockconsumerServiceWriterMockRecorder is the mock recorder for MockconsumerServiceWriter
type MockconsumerServiceWriterMockRecorder struct {
	mock *MockconsumerServiceWriter
}

// NewMockconsumerServiceWriter creates a new mock instance
func NewMockconsumerServiceWriter(ctrl *gomock.Controller) *MockconsumerServiceWriter {
	mock := &MockconsumerServiceWriter{ctrl: ctrl}
	mock.recorder = &MockconsumerServiceWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockconsumerServiceWriter) EXPECT() *MockconsumerServiceWriterMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockconsumerServiceWriter) Write(rm *producer.RefCountedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Write", rm)
}

// Write indicates an expected call of Write
func (mr *MockconsumerServiceWriterMockRecorder) Write(rm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockconsumerServiceWriter)(nil).Write), rm)
}

// Init mocks base method
func (m *MockconsumerServiceWriter) Init(arg0 initType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockconsumerServiceWriterMockRecorder) Init(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockconsumerServiceWriter)(nil).Init), arg0)
}

// Close mocks base method
func (m *MockconsumerServiceWriter) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockconsumerServiceWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockconsumerServiceWriter)(nil).Close))
}

// SetMessageTTLNanos mocks base method
func (m *MockconsumerServiceWriter) SetMessageTTLNanos(value int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMessageTTLNanos", value)
}

// SetMessageTTLNanos indicates an expected call of SetMessageTTLNanos
func (mr *MockconsumerServiceWriterMockRecorder) SetMessageTTLNanos(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMessageTTLNanos", reflect.TypeOf((*MockconsumerServiceWriter)(nil).SetMessageTTLNanos), value)
}

// RegisterFilter mocks base method
func (m *MockconsumerServiceWriter) RegisterFilter(fn producer.FilterFunc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterFilter", fn)
}

// RegisterFilter indicates an expected call of RegisterFilter
func (mr *MockconsumerServiceWriterMockRecorder) RegisterFilter(fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterFilter", reflect.TypeOf((*MockconsumerServiceWriter)(nil).RegisterFilter), fn)
}

// UnregisterFilter mocks base method
func (m *MockconsumerServiceWriter) UnregisterFilter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnregisterFilter")
}

// UnregisterFilter indicates an expected call of UnregisterFilter
func (mr *MockconsumerServiceWriterMockRecorder) UnregisterFilter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterFilter", reflect.TypeOf((*MockconsumerServiceWriter)(nil).UnregisterFilter))
}
