package measurements

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
	atlas "github.com/keltia/ripe-atlas"
)

func TestHandler_CreateMeasurements(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockMeasurementService(ctrl)
	ripeMgr := ripemgr.NewMockRipeMgr(ctrl)

	probes := "111,123"
	miners := []*models.Miner{{
		Address: "fx002",
		IP:      "100.12.35.5",
	}}

	measurements := []*atlas.Measurement{{
		Description: "+++++++++++++++++++++++++",
		FirstHop:    5555,
		Group:       "sdzfsadf",
		GroupID:     8880,
	}}

	service.EXPECT().GetProbIDs().Return(strings.Split(probes, ",")).MaxTimes(1)
	service.EXPECT().GetMiners().Return(miners).MaxTimes(1)
	ripeMgr.EXPECT().CreateMeasurements(miners, probes).Return(measurements, nil).MaxTimes(1)
	service.EXPECT().CreateMeasurements(gomock.Any()).MaxTimes(1)

	h := &Handler{
		Service: service,
		ripeMgr: ripeMgr,
	}

	h.CreateMeasurements()
}

/* return error
CreateMeasurements.DoAndReturn(func(arg0 []*models.Miner, arg1 string) ([]*atlas.Measurement, error) {
		fmt.Print("doo")
		return []*atlas.Measurement{}, fmt.Errorf("---")
	})*/
