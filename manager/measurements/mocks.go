// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ConsenSys/fc-latency-map/manager/measurements (interfaces: MeasurementService)

// Package measurements is a generated GoMock package.
package measurements

import (
	reflect "reflect"

	models "github.com/ConsenSys/fc-latency-map/manager/models"
	gomock "github.com/golang/mock/gomock"
	ripe_atlas "github.com/keltia/ripe-atlas"
)

// MockMeasurementService is a mock of MeasurementService interface.
type MockMeasurementService struct {
	ctrl     *gomock.Controller
	recorder *MockMeasurementServiceMockRecorder
}

// MockMeasurementServiceMockRecorder is the mock recorder for MockMeasurementService.
type MockMeasurementServiceMockRecorder struct {
	mock *MockMeasurementService
}

// NewMockMeasurementService creates a new mock instance.
func NewMockMeasurementService(ctrl *gomock.Controller) *MockMeasurementService {
	mock := &MockMeasurementService{ctrl: ctrl}
	mock.recorder = &MockMeasurementServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMeasurementService) EXPECT() *MockMeasurementServiceMockRecorder {
	return m.recorder
}

// GetLocationsAsPlaces mocks base method.
func (m *MockMeasurementService) GetLocationsAsPlaces() ([]Place, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocationsAsPlaces")
	ret0, _ := ret[0].([]Place)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocationsAsPlaces indicates an expected call of GetLocationsAsPlaces.
func (mr *MockMeasurementServiceMockRecorder) GetLocationsAsPlaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocationsAsPlaces", reflect.TypeOf((*MockMeasurementService)(nil).GetLocationsAsPlaces))
}

// GetMeasurementsRunning mocks base method.
func (m *MockMeasurementService) GetMeasurementsRunning() []*models.Measurement {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMeasurementsRunning")
	ret0, _ := ret[0].([]*models.Measurement)
	return ret0
}

// GetMeasurementsRunning indicates an expected call of GetMeasurementsRunning.
func (mr *MockMeasurementServiceMockRecorder) GetMeasurementsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMeasurementsRunning", reflect.TypeOf((*MockMeasurementService)(nil).GetMeasurementsRunning))
}

// GetMinersWithGeolocation mocks base method.
func (m *MockMeasurementService) GetMinersWithGeolocation() []*models.Miner {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMinersWithGeolocation")
	ret0, _ := ret[0].([]*models.Miner)
	return ret0
}

// GetMinersWithGeolocation indicates an expected call of GetMinersWithGeolocation.
func (mr *MockMeasurementServiceMockRecorder) GetMinersWithGeolocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinersWithGeolocation", reflect.TypeOf((*MockMeasurementService)(nil).GetMinersWithGeolocation))
}

// GetProbIDs mocks base method.
func (m *MockMeasurementService) GetProbIDs(arg0 []Place, arg1, arg2 float64) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProbIDs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetProbIDs indicates an expected call of GetProbIDs.
func (mr *MockMeasurementServiceMockRecorder) GetProbIDs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProbIDs", reflect.TypeOf((*MockMeasurementService)(nil).GetProbIDs), arg0, arg1, arg2)
}

// ImportMeasurement mocks base method.
func (m *MockMeasurementService) ImportMeasurement(arg0 []ripe_atlas.MeasurementResult) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ImportMeasurement", arg0)
}

// ImportMeasurement indicates an expected call of ImportMeasurement.
func (mr *MockMeasurementServiceMockRecorder) ImportMeasurement(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportMeasurement", reflect.TypeOf((*MockMeasurementService)(nil).ImportMeasurement), arg0)
}

// UpsertMeasurements mocks base method.
func (m *MockMeasurementService) UpsertMeasurements(arg0 []*ripe_atlas.Measurement) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpsertMeasurements", arg0)
}

// UpsertMeasurements indicates an expected call of UpsertMeasurements.
func (mr *MockMeasurementServiceMockRecorder) UpsertMeasurements(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMeasurements", reflect.TypeOf((*MockMeasurementService)(nil).UpsertMeasurements), arg0)
}
