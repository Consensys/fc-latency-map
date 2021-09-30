package miners

import (
	"errors"
	"log"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-cid"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/geomgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

var dummyMinerAddress, _ = address.NewFromString("t012345")
var dummyIpAddress = "127.0.0.1"
var dummyMultiAddress, _ = ma.NewMultiaddr("/ip4/" + dummyIpAddress + "/udp/1234")
var dummyMinerInfo = miner.MinerInfo{
	Multiaddrs: [][]byte{dummyMultiAddress.Bytes()},
}
var dummyVerifiedDeals = []fmgr.VerifiedDeal{
	{
		MessageCid: makeCID("dummyCID"),
		Provider:   dummyMinerAddress,
	},
}
var dummyGeoLatitude = 37.39500
var dummyGeoLongitude = -122.08167
var dummyOffset = int(42)
var dummyBlockHeight = int64(42)
var dummyMiner = models.Miner{
	Address:   dummyMinerAddress.String(),
	IP:        dummyIpAddress,
	Latitude:  dummyGeoLatitude,
	Longitude: dummyGeoLongitude,
}

// See https://github.com/filecoin-project/lotus/blob/15d90c24edd3722b71df0b3828667ffde1982d3b/cmd/lotus-health/main_test.go#L160
func makeCID(s string) cid.Cid {
	h1, err := mh.Sum([]byte(s), mh.SHA2_256, -1)
	if err != nil {
		log.Fatal(err)
	}
	return cid.NewCidV1(0x55, h1)
}

func Test_GetAllMiners_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	miners := srv.GetAllMiners()

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_GetAllMiners_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockDbMgr.GetDB().Create(&([]*models.Miner{&dummyMiner}))
	miners := srv.GetAllMiners()

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)

	actual := *(miners[0])
	assert.Equal(t, dummyMiner.Address, actual.Address)
	assert.Equal(t, dummyMiner.IP, actual.IP)
	assert.Equal(t, dummyMiner.Latitude, actual.Latitude)
	assert.Equal(t, dummyMiner.Longitude, actual.Longitude)
}

func Test_ParseMiners_Error_GetBlockHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetBlockHeight().Return(abi.ChainEpoch(int64(0)), errors.New(""))
	miners := srv.ParseMiners(dummyOffset)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMiners_Error_GetVerifiedDealsByBlockRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetBlockHeight().Return(abi.ChainEpoch(dummyBlockHeight), nil)
	mockFMgr.EXPECT().GetVerifiedDealsByBlockRange(gomock.Any(), gomock.Any()).Return(make([]fmgr.VerifiedDeal, 0), errors.New(""))
	miners := srv.ParseMiners(dummyOffset)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMiners_Empty_GetVerifiedDealsByBlockRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetBlockHeight().Return(abi.ChainEpoch(dummyBlockHeight), nil)
	mockFMgr.EXPECT().GetVerifiedDealsByBlockRange(gomock.Any(), gomock.Any()).Return(make([]fmgr.VerifiedDeal, 0), nil)
	miners := srv.ParseMiners(dummyOffset)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMiners_Error_GetMinerInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetBlockHeight().Return(abi.ChainEpoch(dummyBlockHeight), nil)
	mockFMgr.EXPECT().GetVerifiedDealsByBlockRange(gomock.Any(), gomock.Any()).Return(dummyVerifiedDeals, nil)
	mockFMgr.EXPECT().GetMinerInfo(gomock.Any()).Return(miner.MinerInfo{}, errors.New(""))
	miners := srv.ParseMiners(dummyOffset)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMiners_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetBlockHeight().Return(abi.ChainEpoch(dummyBlockHeight), nil)
	mockFMgr.EXPECT().GetVerifiedDealsByBlockRange(gomock.Any(), gomock.Any()).Return(dummyVerifiedDeals, nil)
	mockFMgr.EXPECT().GetMinerInfo(gomock.Any()).Return(dummyMinerInfo, nil)
	mockGMgr.EXPECT().IPGeolocation(gomock.Any()).Return(dummyGeoLatitude, dummyGeoLongitude)
	miners := srv.ParseMiners(dummyOffset)

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)

	actual := *(miners[0])
	assert.Equal(t, dummyMiner.Address, actual.Address)
	assert.Equal(t, dummyMiner.IP, actual.IP)
	assert.Equal(t, dummyMiner.Latitude, actual.Latitude)
	assert.Equal(t, dummyMiner.Longitude, actual.Longitude)
}

func Test_ParseMinersByBlockHeight_Error_GetVerifiedDealsByBlockHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetVerifiedDealsByBlockHeight(gomock.Any()).Return(make([]fmgr.VerifiedDeal, 0), errors.New(""))
	miners := srv.ParseMinersByBlockHeight(dummyBlockHeight)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMinersByBlockHeight_Empty_GetVerifiedDealsByBlockHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetVerifiedDealsByBlockHeight(gomock.Any()).Return(make([]fmgr.VerifiedDeal, 0), nil)
	miners := srv.ParseMinersByBlockHeight(dummyBlockHeight)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMinersByBlockHeight_Error_GetMinerInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetVerifiedDealsByBlockHeight(gomock.Any()).Return(dummyVerifiedDeals, nil)
	mockFMgr.EXPECT().GetMinerInfo(gomock.Any()).Return(miner.MinerInfo{}, errors.New(""))
	miners := srv.ParseMinersByBlockHeight(dummyBlockHeight)

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_ParseMinersByBlockHeight_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockDbMgr := db.NewMockDatabaseMgr()
	mockFMgr := fmgr.NewMockFilecoinMgr(ctrl)
	mockGMgr := geomgr.NewMockGeoMgr(ctrl)
	srv := NewMinerServiceImpl(mockConfig, mockDbMgr, mockFMgr, mockGMgr)

	// Act
	mockFMgr.EXPECT().GetVerifiedDealsByBlockHeight(gomock.Any()).Return(dummyVerifiedDeals, nil)
	mockFMgr.EXPECT().GetMinerInfo(gomock.Any()).Return(dummyMinerInfo, nil)
	mockGMgr.EXPECT().IPGeolocation(gomock.Any()).Return(dummyGeoLatitude, dummyGeoLongitude)
	miners := srv.ParseMinersByBlockHeight(dummyBlockHeight)

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)

	actual := *(miners[0])
	assert.Equal(t, dummyMiner.Address, actual.Address)
	assert.Equal(t, dummyMiner.IP, actual.IP)
	assert.Equal(t, dummyMiner.Latitude, actual.Latitude)
	assert.Equal(t, dummyMiner.Longitude, actual.Longitude)
}
