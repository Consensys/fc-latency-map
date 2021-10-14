package webhook

import (
	"testing"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/stretchr/testify/assert"
)

func Test_Notify_Fail_FilesIsNil(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	notif := NewNotifier(mockConfig)

	// Act
	done := notif.Notify(nil)

	// Assert
	assert.False(t, done)
}

func Test_Notify_Fail_FilesIsEmpty(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	notif := NewNotifier(mockConfig)

	// Act
	done := notif.Notify(&[]string{})

	// Assert
	assert.False(t, done)
}

func Test_Notify_OK_One(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	notif := NewNotifier(mockConfig)

	// Act
	done := notif.Notify(&[]string{"export_2021-10-13.json"})

	// Assert
	assert.True(t, done)
}

func Test_Notify_OK_Many(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	notif := NewNotifier(mockConfig)

	// Act
	done := notif.Notify(&[]string{"export_2021-10-12.json", "export_2021-10-13.json"})

	// Assert
	assert.True(t, done)
}
