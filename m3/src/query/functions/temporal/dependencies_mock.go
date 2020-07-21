// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/query/functions/temporal/dependencies.go

// Copyright (c) 2019 Uber Technologies, Inc.
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

// Package temporal is a generated GoMock package.
package temporal

import (
	"reflect"

	"github.com/m3db/m3/src/query/block"
	"github.com/m3db/m3/src/query/models"

	"github.com/golang/mock/gomock"
)

// Mockcontroller is a mock of controller interface
type Mockcontroller struct {
	ctrl     *gomock.Controller
	recorder *MockcontrollerMockRecorder
}

// MockcontrollerMockRecorder is the mock recorder for Mockcontroller
type MockcontrollerMockRecorder struct {
	mock *Mockcontroller
}

// NewMockcontroller creates a new mock instance
func NewMockcontroller(ctrl *gomock.Controller) *Mockcontroller {
	mock := &Mockcontroller{ctrl: ctrl}
	mock.recorder = &MockcontrollerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockcontroller) EXPECT() *MockcontrollerMockRecorder {
	return m.recorder
}

// BlockBuilder mocks base method
func (m *Mockcontroller) BlockBuilder(queryCtx *models.QueryContext, blockMeta block.Metadata, seriesMeta []block.SeriesMeta) (block.Builder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockBuilder", queryCtx, blockMeta, seriesMeta)
	ret0, _ := ret[0].(block.Builder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockBuilder indicates an expected call of BlockBuilder
func (mr *MockcontrollerMockRecorder) BlockBuilder(queryCtx, blockMeta, seriesMeta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockBuilder", reflect.TypeOf((*Mockcontroller)(nil).BlockBuilder), queryCtx, blockMeta, seriesMeta)
}

// Process mocks base method
func (m *Mockcontroller) Process(queryCtx *models.QueryContext, block block.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", queryCtx, block)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process
func (mr *MockcontrollerMockRecorder) Process(queryCtx, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*Mockcontroller)(nil).Process), queryCtx, block)
}
