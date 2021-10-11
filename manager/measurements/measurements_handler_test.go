package measurements

import (
	"fmt"
	"github.com/golang/mock/gomock"
	atlas "github.com/keltia/ripe-atlas"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

func TestHandler_CreateMeasurementsRipeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	probes := "111,123"
	miners := []*models.Miner{{
		Address: "fx002",
		IP:      "100.12.35.5",
	}}
	service.EXPECT().getLocationsAsPlaces().Return([]Place{}, nil)
	service.EXPECT().getProbIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(strings.Split(probes, ",")).MaxTimes(1)
	service.EXPECT().GetMinersWithGeolocation().Return(miners).MaxTimes(1)

	ripeMgr.EXPECT().
		CreateMeasurements(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error")).
		MaxTimes(1)

	service.EXPECT().UpsertMeasurements(gomock.Any()).Times(0)

	h.createMeasurements(nil)
}

func TestHandler_CreateMeasurements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	probes := "111,123"
	miners := []*models.Miner{{
		Address: "fx002",
		IP:      "100.12.35.5",
	}}

	service.EXPECT().getLocationsAsPlaces().Return([]Place{}, nil)
	service.EXPECT().getProbIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(strings.Split(probes, ",")).Times(1)

	service.EXPECT().GetMinersWithGeolocation().Return(miners).Times(1)

	ripeMgr.EXPECT().
		CreateMeasurements(miners, probes, 0).
		Return(nil, nil).
		MaxTimes(1)

	service.EXPECT().UpsertMeasurements(gomock.Any()).Times(1)

	h.createMeasurements(nil)
}

func TestHandler_ImportMeasuresError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	service.EXPECT().getMeasurementsRunning().Times(1)
	service.EXPECT().ImportMeasurement(gomock.Any()).Times(0)

	h.importMeasures()
}

func TestHandler_ImportMeasures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMeas := []atlas.Measurement{
		{Af: 4, ID: 1},
	}

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	service.EXPECT().
		getMeasurementsRunning().
		Return([]*models.Measurement{
			{
				Model:          gorm.Model{ID: 0},
				IsOneOff:       false,
				MeasurementID:  1,
				StartTime:      0,
				StopTime:       999999999,
				Status:         "runnig",
				StatusStopTime: 0,
			},
		}).
		Times(1)

	ripeMgr.EXPECT().GetMeasurement(1).Return(&mockMeas[0], nil)
	service.EXPECT().
		UpsertMeasurements(gomock.Any()).
		Times(1)

	ripeMgr.EXPECT().GetMeasurementResults(1).
		Return(
			[]atlas.MeasurementResult{
				{Af: 4, Avg: 1, DstAddr: "1.1.1.1", PrbID: 1, MsmID: 1},
			}, nil)

	service.EXPECT().ImportMeasurement(gomock.Any()).Times(1)

	h.importMeasures()
}

func TestHandler_ImportMeasures_GetMeasurementError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMeas := []atlas.Measurement{
		{Af: 4, ID: 1},
	}

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	service.EXPECT().
		getMeasurementsRunning().
		Return([]*models.Measurement{
			{
				Model:          gorm.Model{ID: 0},
				IsOneOff:       false,
				MeasurementID:  1,
				StartTime:      0,
				StopTime:       999999999,
				Status:         "runnig",
				StatusStopTime: 0,
			},
		}).
		Times(1)

	ripeMgr.EXPECT().GetMeasurement(1).Return(&mockMeas[0], errors.New("error"))

	service.EXPECT().
		UpsertMeasurements(gomock.Any()).
		Times(0)

	ripeMgr.EXPECT().GetMeasurementResults(1).
		Return(
			[]atlas.MeasurementResult{
				{Af: 4, Avg: 1, DstAddr: "1.1.1.1", PrbID: 1, MsmID: 1},
			}, nil).Times(0)

	service.EXPECT().ImportMeasurement(gomock.Any()).Times(0)

	h.importMeasures()
}

func TestHandler_ImportMeasures_GetMeasurementResultsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMeas := []atlas.Measurement{
		{Af: 4, ID: 1},
	}

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	service.EXPECT().
		getMeasurementsRunning().
		Return([]*models.Measurement{
			{
				Model:          gorm.Model{ID: 0},
				IsOneOff:       false,
				MeasurementID:  1,
				StartTime:      0,
				StopTime:       999999999,
				Status:         "runnig",
				StatusStopTime: 0,
			},
		}).
		Times(1)

	ripeMgr.EXPECT().GetMeasurement(1).Return(&mockMeas[0], nil)

	service.EXPECT().
		UpsertMeasurements(gomock.Any()).
		Times(1)

	ripeMgr.EXPECT().GetMeasurementResults(1).
		Return(
			[]atlas.MeasurementResult{
				{Af: 4, Avg: 1, DstAddr: "1.1.1.1", PrbID: 1, MsmID: 1},
			}, errors.New("error"))

	service.EXPECT().ImportMeasurement(gomock.Any()).Times(0)

	h.importMeasures()
}