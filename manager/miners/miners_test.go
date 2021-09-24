package miners

import (
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
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

func Test_ParseMiners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyOffset := uint(42)
	dummyBlockHeight := int64(42)
	dummyVerifiedDeals := make([]fmgr.VerifiedDeal, 0)
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	mockFMgr.EXPECT().GetBlockHeight().Return(abi.ChainEpoch(dummyBlockHeight), nil)
	mockFMgr.EXPECT().GetVerifiedDealsByBlockRange(abi.ChainEpoch(dummyBlockHeight), dummyOffset).Return(dummyVerifiedDeals, nil)

	miners := srv.ParseMiners(dummyOffset)
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMinersByBlockHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyBlockHeight := int64(42)
	dummyVerifiedDeals := make([]fmgr.VerifiedDeal, 0)
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	mockFMgr.EXPECT().GetVerifiedDealsByBlockHeight(abi.ChainEpoch(dummyBlockHeight)).Return(dummyVerifiedDeals, nil)

	miners := srv.ParseMinersByBlockHeight(dummyBlockHeight)
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}
