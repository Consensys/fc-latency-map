package webhook

import (
	"strings"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

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
	defer gock.Off() // Flush pending mocks after test execution

	// Arrange
	mockConfig := config.NewMockConfig()
	notif := NewNotifier(mockConfig)
	urls := strings.Split(mockConfig.GetString("WEBHOOK_NOTIFY_URLS"), ",")
	for _, url := range urls {
		gock.New(url).
			Reply(200).
			JSON(map[string]string{"status": "success"})
	}

	// Act
	done := notif.Notify(&[]string{"export_2021-10-13.json"})

	// Assert
	assert.True(t, done)
}

func Test_Notify_OK_Many(t *testing.T) {
	// Arrange
	mockConfig := config.NewMockConfig()
	notif := NewNotifier(mockConfig)
	urls := strings.Split(mockConfig.GetString("WEBHOOK_NOTIFY_URLS"), ",")
	for _, url := range urls {
		gock.New(url).
			Reply(200).
			JSON(map[string]string{"status": "success"})
	}

	// Act
	done := notif.Notify(&[]string{"export_2021-10-12.json", "export_2021-10-13.json"})

	// Assert
	assert.True(t, done)
}
