// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ConsenSys/fc-latency-map/manager/ripemgr (interfaces: RipeMgr)

// Package ripemgr is a generated GoMock package.
package ripemgr

import (
	reflect "reflect"

	models "github.com/ConsenSys/fc-latency-map/manager/models"
	gomock "github.com/golang/mock/gomock"
	atlas "github.com/keltia/ripe-atlas"
)

// MockRipeMgr is a mock of RipeMgr interface.
type MockRipeMgr struct {
	ctrl     *gomock.Controller
	recorder *MockRipeMgrMockRecorder
}

// MockRipeMgrMockRecorder is the mock recorder for MockRipeMgr.
type MockRipeMgrMockRecorder struct {
	mock *MockRipeMgr
}

// NewMockRipeMgr creates a new mock instance.
func NewMockRipeMgr(ctrl *gomock.Controller) *MockRipeMgr {
	mock := &MockRipeMgr{ctrl: ctrl}
	mock.recorder = &MockRipeMgrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRipeMgr) EXPECT() *MockRipeMgrMockRecorder {
	return m.recorder
}

// CreateMeasurements mocks base method.
func (m *MockRipeMgr) CreateMeasurements(arg0 []*models.Miner, arg1 string, arg2 int) ([]*atlas.Measurement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMeasurements", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*atlas.Measurement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMeasurements indicates an expected call of CreateMeasurements.
func (mr *MockRipeMgrMockRecorder) CreateMeasurements(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMeasurements", reflect.TypeOf((*MockRipeMgr)(nil).CreateMeasurements), arg0, arg1, arg2)
}

// GetMeasurement mocks base method.
func (m *MockRipeMgr) GetMeasurement(arg0 int) (*atlas.Measurement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMeasurement", arg0)
	ret0, _ := ret[0].(*atlas.Measurement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMeasurement indicates an expected call of GetMeasurement.
func (mr *MockRipeMgrMockRecorder) GetMeasurement(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMeasurement", reflect.TypeOf((*MockRipeMgr)(nil).GetMeasurement), arg0)
}

// GetMeasurementResults mocks base method.
func (m *MockRipeMgr) GetMeasurementResults(arg0, arg1 int) ([]atlas.MeasurementResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMeasurementResults", arg0, arg1)
	ret0, _ := ret[0].([]atlas.MeasurementResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMeasurementResults indicates an expected call of GetMeasurementResults.
func (mr *MockRipeMgrMockRecorder) GetMeasurementResults(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMeasurementResults", reflect.TypeOf((*MockRipeMgr)(nil).GetMeasurementResults), arg0, arg1)
}

// GetProbes mocks base method.
func (m *MockRipeMgr) GetProbes(arg0 map[string]string) ([]atlas.Probe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProbes", arg0)
	ret0, _ := ret[0].([]atlas.Probe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProbes indicates an expected call of GetProbes.
func (mr *MockRipeMgrMockRecorder) GetProbes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProbes", reflect.TypeOf((*MockRipeMgr)(nil).GetProbes), arg0)
}
