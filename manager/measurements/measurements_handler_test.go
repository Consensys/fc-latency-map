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
	service.EXPECT().GetProbIDs().Return(strings.Split(probes, ",")).MaxTimes(1)
	service.EXPECT().GetMiners().Return(miners).MaxTimes(1)

	ripeMgr.EXPECT().
		CreateMeasurements(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error")).
		MaxTimes(1)

	service.EXPECT().CreateMeasurements(gomock.Any()).Times(0)

	h.CreateMeasurements()
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

	service.EXPECT().GetProbIDs().Return(strings.Split(probes, ",")).Times(1)

	service.EXPECT().GetMiners().Return(miners).Times(1)

	ripeMgr.EXPECT().
		CreateMeasurements(miners, probes).
		Return(nil, nil).
		MaxTimes(1)

	service.EXPECT().CreateMeasurements(gomock.Any()).Times(1)

	h.CreateMeasurements()
}

func TestHandler_GetMeasuresRipeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}
	service.EXPECT().GetMeasuresLastResultTime().Times(1)
	ripeMgr.EXPECT().
		GetMeasurementResults(gomock.Any()).
		Return(nil, fmt.Errorf("error")).
		MaxTimes(1)

	service.EXPECT().ImportMeasurement(gomock.Any()).Times(0)

	h.GetMeasures()
}
func TestHandler_GetMeasures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}
	service.EXPECT().GetMeasuresLastResultTime().Times(1)

	ripeMgr.EXPECT().
		GetMeasurementResults(gomock.Any()).
		Return(nil, nil).
		MaxTimes(1)

	service.EXPECT().ImportMeasurement(gomock.Any()).Times(1)

	h.GetMeasures()
}
