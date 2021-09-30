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
	locsList := []*models.Location{}
	srv.DBMgr.GetDB().Find(&locsList)
	var bestProbes []*atlas.Probe

	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
			"iata":    location.IataCode,
		}).Info("Get probes for airport")

		nearestProbe, err := srv.RipeMgr.GetNearestProbe(location.GeoLocation.Latitude, location.GeoLocation.Longitude)
		if err != nil {
			return nil, err
		}

		bestProbes = append(bestProbes, nearestProbe)
	}

	return bestProbes, nil
}

func (srv *ProbeServiceImpl) GetAllProbes() []*models.Probe {
	probesList := []*models.Probe{}
	srv.DBMgr.GetDB().Find(&probesList)
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
			newProbe.GeoLocation = models.GeoLocation{
				Latitude:  probe.Geometry.Coordinates[0],
				Longitude: probe.Geometry.Coordinates[1],
			}
		}

		probeExits := models.Probe{}
		srv.DBMgr.GetDB().Where("probe_id = ?", probe.ID).First(&probeExits)

		if (models.Probe{}) == probeExits {
			err := srv.DBMgr.GetDB().Debug().Model(&models.Probe{}).Create(&newProbe).Error
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("unable to insert probes")
			}
			log.Printf("Add new location, ID: %v", newProbe.ProbeID)
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
			Country: probe.CountryCode,
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
