package export

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var dummyFile = &[]string{
	"export_20211020.json",
	"export_20211019.json",
}

func Test_Export_OK_Nil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockSrv := NewMockService(ctrl)
	hdlr := &ExportHandler{
		Service: mockSrv,
	}
	mockSrv.EXPECT().export().Return(nil)

	// Act
	files := hdlr.Export()

	// Assert
	assert.Nil(t, files)
}

func Test_Export_OK_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockSrv := NewMockService(ctrl)
	hdlr := &ExportHandler{
		Service: mockSrv,
	}
	mockSrv.EXPECT().export().Return(&[]string{})

	// Act
	files := hdlr.Export()

	// Assert
	assert.NotNil(t, files)
	assert.Empty(t, *files)
}

func Test_Export_OK_NotEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockSrv := NewMockService(ctrl)
	hdlr := &ExportHandler{
		Service: mockSrv,
	}
	mockSrv.EXPECT().export().Return(dummyFile)

	// Act
	files := hdlr.Export()

	// Assert
	assert.NotNil(t, files)
	assert.Equal(t, *dummyFile, *files)
}
