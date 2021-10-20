package ripemgr

import (
	"fmt"

	"testing"

	"github.com/golang/mock/gomock"
	atlas "github.com/keltia/ripe-atlas"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

func TestNewRipeImpl2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIndex := NewMockRipeMgr(ctrl)
	b := "111"
	a := []*models.Miner{{
		Address: "",
		IP:      "",
	}}

	createMeasurements := mockIndex.EXPECT().CreateMeasurements(a, "111", 0)
	createMeasurements.DoAndReturn(func(_ []*models.Miner, _ string, _ int) ([]*atlas.Measurement, error) {
		fmt.Print("doo")
		return []*atlas.Measurement{}, fmt.Errorf("error")
	})
	createMeasurements.MaxTimes(5)
	measurements, err := mockIndex.CreateMeasurements(a, b, 0)
	assert.NotNil(t, err)
	assert.NotNil(t, measurements)
	createMeasurements.Return([]*atlas.Measurement{{
		Description: "+++++++++++++++++++++++++",
		FirstHop:    5555,
		Group:       "sdzfsadf",
		GroupID:     8880,
	}}, nil)

	i, err := mockIndex.CreateMeasurements(a, b, 0)
	assert.Nil(t, err)
	assert.Equal(t, []*atlas.Measurement{{
		Description: "+++++++++++++++++++++++++",
		FirstHop:    5555,
		Group:       "sdzfsadf",
		GroupID:     8880,
	}}, i)
}
