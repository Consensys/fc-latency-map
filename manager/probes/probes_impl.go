package probes

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"

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

func (srv *ProbeServiceImpl) RequestProbes() error {
	dbc := srv.DBMgr.GetDB()
	var places []measurements.Place
	err := srv.DBMgr.GetDB().Model(&models.Probe{}).Find(&places).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get places from db")

		return err
	}

	locsList := []*models.Location{}
	dbc.
		Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}, Desc: true}).
		Order("country, iata_code").
		Find(&locsList)
	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
			"iata":    location.IataCode,
		}).Info("Get probes for airport")

		nearestProbeIDs := measurements.FindNearest(places,
			measurements.Place{Latitude: location.Latitude, Longitude: location.Longitude},
			config.NewConfig().GetInt("RIPE_PROBES_PER_AIRPORT"))

		var probes []*models.Probe
		err := dbc.Model(&models.Probe{}).Find(&probes, "id in ?", nearestProbeIDs).Error
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("get places from db")

			return err
		}

		location.Probes = probes
		dbc.Updates(location)
	}

	return nil
}

func (srv *ProbeServiceImpl) ListProbes() []*models.Probe {
	probesList := []*models.Probe{}
	srv.DBMgr.GetDB().Find(&probesList)
	for _, probe := range probesList {
		log.Printf("Probe ID: %d - Country code: %s \n", probe.ProbeID, probe.CountryCode)
	}

	return probesList
}

func (srv *ProbeServiceImpl) GetTotalProbes() int64 {
	var count int64
	srv.DBMgr.GetDB().Model(&models.Probe{}).Count(&count)
	return count
}

func (srv *ProbeServiceImpl) Update() {
	err := srv.RequestProbes()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("List")

		return
	}

	log.Println("Probes successfully updated")
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
	for _, v := range probes {
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
