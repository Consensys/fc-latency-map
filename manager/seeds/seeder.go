package seeds

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	"gorm.io/gorm"
)

// Execute runs the data seed process
func Execute(db *gorm.DB) error {
	for i, _ := range locations {
		err := db.Debug().Model(&models.Location{}).Create(&locations[i]).Error
		if err != nil {
			return err
		}
		log.Printf("Location ID: %v", locations[i].ID)
	}
	return nil
}