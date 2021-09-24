package miners

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/geomgr"
	gomock "github.com/golang/mock/gomock"
)

func Test_GetAllMiners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	miners := srv.GetAllMiners()
	assert.NotNil(t, miners)
}
