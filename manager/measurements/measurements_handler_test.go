package measurements

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

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
	service.EXPECT().GetLocationsAsPlaces().Return([]Place{}, nil)
	service.EXPECT().GetProbIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(strings.Split(probes, ",")).MaxTimes(1)
	service.EXPECT().GetMinersWithGeolocation().Return(miners).MaxTimes(1)

	ripeMgr.EXPECT().
		CreateMeasurements(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error")).
		MaxTimes(1)

	service.EXPECT().UpsertMeasurements(gomock.Any()).Times(0)

	h.CreateMeasurements(nil)
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

	service.EXPECT().GetLocationsAsPlaces().Return([]Place{}, nil)
	service.EXPECT().GetProbIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(strings.Split(probes, ",")).Times(1)

	service.EXPECT().GetMinersWithGeolocation().Return(miners).Times(1)

	ripeMgr.EXPECT().
		CreateMeasurements(miners, probes, 0).
		Return(nil, nil).
		MaxTimes(1)

	service.EXPECT().UpsertMeasurements(gomock.Any()).Times(1)

	h.CreateMeasurements(nil)
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

	service.EXPECT().GetMeasurementsRunning().Times(1)
	service.EXPECT().ImportMeasurement(gomock.Any()).Times(0)

	h.ImportMeasures()
}

func TestHandler_ImportMeasures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	service.EXPECT().GetMeasurementsRunning().Times(1)

	ripeMgr.EXPECT().
		GetMeasurementResults(gomock.Any()).
		Return(nil, nil).
		MaxTimes(1)

	h.ImportMeasures()
}
