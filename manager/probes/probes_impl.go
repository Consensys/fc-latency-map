package probes

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type ProbeServiceImpl struct {
	DBMgr   db.DatabaseMgr
	RipeMgr ripemgr.RipeMgr
}

const point = "Point"

func NewProbeServiceImpl(dbMgr db.DatabaseMgr, ripeMgr ripemgr.RipeMgr) (ProbeService, error) {
	return &ProbeServiceImpl{
		DBMgr:   dbMgr,
		RipeMgr: ripeMgr,
	}, nil
}

func (srv *ProbeServiceImpl) RequestProbes() ([]*models.Probe, error) {
	locsList := []*models.Location{}
	srv.DBMgr.GetDB().Order("country, iata_code").Find(&locsList)
	var bestProbes []*models.Probe

	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
			"iata":    location.IataCode,
		}).Info("Get probes for airport")

		nearestProbe, err := srv.RipeMgr.GetNearestProbe(location.Latitude, location.Longitude)
		if err != nil {
			return nil, err
		}

		for _, v := range *nearestProbe { // nolint:gocritic
			newProbe := &models.Probe{
				ProbeID:     v.ID,
				CountryCode: v.CountryCode,
				Location:    *location,
			}
			if v.Geometry.Type == point {
				newProbe.Latitude = v.Geometry.Coordinates[0]
				newProbe.Longitude = v.Geometry.Coordinates[1]
			}

			bestProbes = append(bestProbes, newProbe)
		}
	}

	return bestProbes, nil
}

func (srv *ProbeServiceImpl) ListProbes() []*models.Probe {
	probesList := []*models.Probe{}
	srv.DBMgr.GetDB().Find(&probesList)
	for _, probe := range probesList {
		log.Printf("Probe ID: %d - Country code: %s - IataCode: %s\n", probe.ProbeID, probe.CountryCode, probe.IataCode)
	}

	return probesList
}

func (srv *ProbeServiceImpl) Update() {
	probes, err := srv.RequestProbes()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("List")

		return
	}

	// update with new probes
	for _, probe := range probes {
		probeExits := models.Probe{}
		srv.DBMgr.GetDB().Where("probe_id = ?", probe.ID).First(&probeExits)

		if (models.Probe{}) == probeExits {
			err := srv.DBMgr.GetDB().Model(&models.Probe{}).
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "probe_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"latitude", "longitude"}),
				}).Create(probe).Error
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("unable to insert probes")
			}
			log.Printf("Add new probr, ID: %v", probe.ProbeID)
		} else {
			log.Printf("Probe already exists, Probe ID: %v", probeExits.ProbeID)
		}
	}

	srv.removeDeprecated()

	log.Println("Probes successfully updated")
}

// removeDeprecated update by removing probes not in location list
func (srv *ProbeServiceImpl) removeDeprecated() {
	probesList := []*models.Probe{}
	srv.DBMgr.GetDB().Find(&probesList)
	for _, probe := range probesList {
		location := models.Location{
			IataCode: probe.IataCode,
		}
		locationExists := models.Location{}
		srv.DBMgr.GetDB().Where(&location).First(&locationExists)
		if locationExists == (models.Location{}) {
			srv.DBMgr.GetDB().Delete(&models.Probe{}, probe.ID)
			log.WithFields(log.Fields{
				"ID": probe.ID,
			}).Info("Probe deleted")
		}
	}
}
func (srv *ProbeServiceImpl) ImportProbes() {
	opts := make(map[string]string)
	opts["status_name"] = "Connected"
	opts["sort"] = "id"

	probes, err := srv.RipeMgr.GetProbes(opts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("get all probe from ripe")
		return
	}

	probesDB := []*models.Probe{}
	for _, v := range probes { // nolint:gocritic
		newProbe := &models.Probe{
			ProbeID:     v.ID,
			CountryCode: v.CountryCode,
			Status:      v.Status.Name,
		}
		if v.Geometry.Type == point {
			newProbe.Latitude = v.Geometry.Coordinates[0]
			newProbe.Longitude = v.Geometry.Coordinates[1]
		}
		probesDB = append(probesDB, newProbe)
	}
	dbc := srv.DBMgr.GetDB()

	// update all the rows
	err = dbc.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Model(&models.Probe{}).
		Updates(models.Probe{Status: "Disconnected"}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to delete deprecated probes")
	}

	// upsert to new values from ripe
	err = dbc.Model(&models.Probe{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "probe_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"latitude", "longitude", "status"}),
		}).Create(probesDB).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to insert probes")
	}
}
