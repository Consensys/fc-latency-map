package miners

import (
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var dummyMiners = []*models.Miner{
	&models.Miner{
		Address:   "f01234",
		IP:        "127.0.0.1",
		Latitude:  float64(37.6597400),
		Longitude: float64(-97.5753300),
		Port:      int(80),
	},
}

func Test_OK_GetAllMiners_Nil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := &MinerHandler{
		Conf: *mockConfig,
		MSer: mockMinerSrv,
	}
	mockMinerSrv.EXPECT().GetAllMiners().Return(nil)

	// Act
	miners := hdlr.GetAllMiners()

	// Assert
	assert.Nil(t, miners)
}

func Test_OK_GetAllMiners_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := NewMinerHandler(mockConfig, mockMinerSrv)
	mockMinerSrv.EXPECT().GetAllMiners().Return([]*models.Miner{})

	// Act
	miners := hdlr.GetAllMiners()

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_OK_GetAllMiners_NotEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := NewMinerHandler(mockConfig, mockMinerSrv)
	mockMinerSrv.EXPECT().GetAllMiners().Return(dummyMiners)

	// Act
	miners := hdlr.GetAllMiners()

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)
	assert.Equal(t, *dummyMiners[0], *miners[0])
}

func Test_OK_MinersParseOffset_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := NewMinerHandler(mockConfig, mockMinerSrv)
	mockMinerSrv.EXPECT().ParseMinersByBlockOffset(gomock.Any()).Return(dummyMiners)

	// Act
	miners := hdlr.MinersParseOffset("")

	// Assert
	assert.NotNil(t, miners)
	assert.Empty(t, miners)
}

func Test_OK_MinersParseOffset_NotEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := NewMinerHandler(mockConfig, mockMinerSrv)
	mockMinerSrv.EXPECT().ParseMinersByBlockOffset(gomock.Any()).Return(dummyMiners)

	// Act
	miners := hdlr.MinersParseOffset("42")

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)
	assert.Equal(t, *dummyMiners[0], *miners[0])
}

func Test_OK_MinersParseBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := NewMinerHandler(mockConfig, mockMinerSrv)
	mockMinerSrv.EXPECT().ParseMinersByBlockHeight(gomock.Any()).Return(dummyMiners)

	// Act
	miners := hdlr.MinersParseBlock(42)

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)
	assert.Equal(t, *dummyMiners[0], *miners[0])
}

func Test_OK_MinersParseStateMarket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockConfig := config.NewMockConfig()
	mockMinerSrv := NewMockMinerService(ctrl)
	hdlr := NewMinerHandler(mockConfig, mockMinerSrv)
	mockMinerSrv.EXPECT().ParseMinersByStateMarket().Return(dummyMiners)

	// Act
	miners := hdlr.MinersParseStateMarket()

	// Assert
	assert.NotNil(t, miners)
	assert.NotEmpty(t, miners)
	assert.Equal(t, *dummyMiners[0], *miners[0])
}
