package probes

import (
	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type ProbeServiceImpl struct {
	DbMgr 			*db.DatabaseMgr
	RipeMgr 		*ripemgr.RipeMgr
}

func NewProbeServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr, ripeMgr *ripemgr.RipeMgr) (ProbeService, error) {
	return &ProbeServiceImpl{
		DbMgr: dbMgr,
		RipeMgr:  ripeMgr,
	}, nil
}

func (srv *ProbeServiceImpl) GetBestProbes(countryProbes []atlas.Probe) (atlas.Probe, error) {
	for _, probe := range countryProbes {
		if probe.Status.Name == "Connected" {
			log.WithFields(log.Fields{
				"ID": probe.ID,
			}).Info("Best probe found")
			return probe, nil
		}
	}
	return atlas.Probe{}, nil
}

func (srv *ProbeServiceImpl) RequestProbes() ([]atlas.Probe, error) {
	var locsList = []*models.Location{}
	(*srv.DbMgr).GetDb().Find(&locsList)
	var bestProbes []atlas.Probe

	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
		}).Info("Get probes for country")
		opts := make(map[string]string)
		opts["country_code"] = location.Country
		countryProbes, err := (*srv.RipeMgr).GetProbes(opts)
		if err != nil {
			return nil, err
		}
		bestProbe, err := srv.GetBestProbes(countryProbes)
		if err != nil {
			return nil, err
		}
		bestProbes = append(bestProbes, bestProbe)
	}
	return bestProbes, nil
}

func (srv *ProbeServiceImpl) GetAllProbes() []*models.Probe {
	var probesList = []*models.Probe{}
	(*srv.DbMgr).GetDb().Find(&probesList)
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
			ProbeID: probe.ID,
			CountryCode: probe.CountryCode,
		}

		var probe = models.Probe{}
		(*srv.DbMgr).GetDb().Where("country_code = ?", probe.CountryCode).First(&probe)
		if (models.Probe{}) == probe {
			err := (*srv.DbMgr).GetDb().Debug().Model(&models.Probe{}).Create(&newProbe).Error
			if err != nil {
				panic("Unable to create probe")
			}
			log.Printf("Add new location, ID: %v", newProbe.ProbeID)
		} else {
			log.Printf("Probe already exists, Probe ID: %v", probe.ProbeID)
		}
	}

	// update by removing probes not in location list
	var probesList = []*models.Probe{}
	(*srv.DbMgr).GetDb().Find(&probesList)
	for _, probe := range probesList {
		var location = models.Location{
			Country: probe.CountryCode,
		}
		var locationExists = models.Location{}
		(*srv.DbMgr).GetDb().Where(&location).First(&locationExists)
		if locationExists == (models.Location{}) {
			(*srv.DbMgr).GetDb().Delete(&models.Probe{}, probe.ID)
			log.WithFields(log.Fields{
				"ID": probe.ID,
			}).Info("Probe deleted")
			
		}
	}

	log.Println("Probes successfully updated")	
}
