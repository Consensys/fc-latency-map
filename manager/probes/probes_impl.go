package probes

import (
	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type ProbeServiceImpl struct {
	c         *atlas.Client
	DbMgr 		*db.DatabaseMgr
}

func NewProbeServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr) (ProbeService, error) {
	cfgs := []atlas.Config{}
	cfgs = append(cfgs, atlas.Config{
		APIKey: conf.GetString("RIPE_API_KEY"),
	})
	c, err := atlas.NewClient(cfgs...)
	if err != nil {
		log.Println("Connecting to Ripe Atlas API", err)
		return nil, err
	}
	ver := atlas.GetVersion()
	log.Println("api version ", ver)

	return &ProbeServiceImpl{
		c:  c,
		DbMgr: dbMgr,
	}, nil
}

func (srv *ProbeServiceImpl) GetProbe(id int) (m *atlas.Probe, err error) {
	return srv.c.GetProbe(id)
}

func (srv *ProbeServiceImpl) GetProbes(countryCode string) ([]atlas.Probe, error) {
	opts := make(map[string]string)
	opts["country_code"] = countryCode

	probes, err := srv.c.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil
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

func (srv *ProbeServiceImpl) RequestAllProbes() ([]atlas.Probe, error) {
	var locsList = []*models.Location{}
	(*srv.DbMgr).GetDb().Find(&locsList)
	var bestProbes []atlas.Probe

	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
		}).Info("Get probes for country")
		countryProbes, err := srv.GetProbes(location.Country)
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
	// get countries from db
	probes, err := srv.RequestAllProbes() // by countries
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GetAllProbes")
		return
	}

	// update db probes
	for _, probe := range probes {
		newProbe := models.Probe{
			ProbeID: probe.ID,
			CountryCode: probe.CountryCode,
		}

		var probe = models.Probe{}
		(*srv.DbMgr).GetDb().Where(&newProbe).First(&probe)
		if (models.Probe{}) == probe {
			err := (*srv.DbMgr).GetDb().Debug().Model(&models.Probe{}).Create(&newProbe).Error
			if err != nil {
				panic("Unable to create probe")
			}
			log.Printf("Add new location, ID: %v", newProbe.ProbeID)
		} else {
			log.Printf("Probe already exists, Probe ID: %v", probe.ProbeID)
		}

		(*srv.DbMgr).GetDb().Create(&newProbe)
	}
	log.Println("Probes successfully updated")
	
}
