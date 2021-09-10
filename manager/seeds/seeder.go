package seeds

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/models"
	"gorm.io/gorm"
)

// Execute runs the data seed process
func Execute(db *gorm.DB) error {
	for i, _ := range locations {
		var location = models.Location{}
		db.Where(&locations[i]).First(&location)
		if (models.Location{}) != location  {
			err := db.Debug().Model(&models.Location{}).Create(&locations[i]).Error
			if err != nil {
				return err
			}
			log.Printf("Add new location, ID: %v", locations[i].ID)
		}
		
	}
	return nil
}