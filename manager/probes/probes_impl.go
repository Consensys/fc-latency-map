package probes

import (
	log "github.com/sirupsen/logrus"

	atlas "github.com/keltia/ripe-atlas"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type ProbeServiceImpl struct {
	DBMgr   db.DatabaseMgr
	RipeMgr ripemgr.RipeMgr
}

func NewProbeServiceImpl(dbMgr db.DatabaseMgr, ripeMgr ripemgr.RipeMgr) (ProbeService, error) {
	return &ProbeServiceImpl{
		DBMgr:   dbMgr,
		RipeMgr: ripeMgr,
	}, nil
}

func (srv *ProbeServiceImpl) RequestProbes() ([]*atlas.Probe, error) {
	var locsList = []*models.Location{}
	(srv.DBMgr).GetDB().Find(&locsList)
	var bestProbes []*atlas.Probe

	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
		}).Info("Get probes for country")

		nearestProbe, err := srv.RipeMgr.GetNearestProbe(location.Latitude, location.Longitude)
		if err != nil {
			return nil, err
		}

		bestProbes = append(bestProbes, nearestProbe)
	}

	return bestProbes, nil
}

func (srv *ProbeServiceImpl) GetAllProbes() []*models.Probe {
	var probesList = []*models.Probe{}
	(srv.DBMgr).GetDB().Find(&probesList)
	for _, probe := range probesList {
		log.Printf("Probe ID: %d - Country code: %s\n", probe.ProbeID, probe.CountryCode)
	}
	return probesList
}

func (srv *ProbeServiceImpl) Update() {
	probes, err := srv.RequestProbes()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetAllProbes")
		return
	}

	// update with new probes
	for _, probe := range probes {
		newProbe := models.Probe{
			ProbeID:     probe.ID,
			CountryCode: probe.CountryCode,
		}
		if probe.Geometry.Type == "Point" {
			newProbe.Latitude = probe.Geometry.Coordinates[0]
			newProbe.Longitude = probe.Geometry.Coordinates[1]
		}

		var probeExits = models.Probe{}
		(srv.DBMgr).GetDB().Where("probe_id = ?", probe.ID).First(&probeExits)

		if (models.Probe{}) == probeExits {
			err := (srv.DBMgr).GetDB().Debug().Model(&models.Probe{}).Create(&newProbe).Error
			if err != nil {
				panic("Unable to create probe")
			}
			log.Printf("Add new location, ID: %v", newProbe.ProbeID)
		} else {
			log.Printf("Probe already exists, Probe ID: %v", probeExits.ProbeID)
		}
	}

	// update by removing probes not in location list
	var probesList = []*models.Probe{}
	(srv.DBMgr).GetDB().Find(&probesList)
	for _, probe := range probesList {
		var location = models.Location{
			Country: probe.CountryCode,
		}
		var locationExists = models.Location{}
		(srv.DBMgr).GetDB().Where(&location).First(&locationExists)
		if locationExists == (models.Location{}) {
			(srv.DBMgr).GetDB().Delete(&models.Probe{}, probe.ID)
			log.WithFields(log.Fields{
				"ID": probe.ID,
			}).Info("Probe deleted")
		}
	}

	log.Println("Probes successfully updated")
}
