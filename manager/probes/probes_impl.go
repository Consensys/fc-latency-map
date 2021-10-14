package probes

import (
	atlas "github.com/keltia/ripe-atlas"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/geomgr"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/measurements"

	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/ConsenSys/fc-latency-map/manager/ripemgr"
)

type ProbeServiceImpl struct {
	DBMgr   db.DatabaseMgr
	RipeMgr ripemgr.RipeMgr
	GeoMgr  geomgr.GeoMgr
}

const point = "Point"

func NewProbeServiceImpl(dbMgr db.DatabaseMgr, ripeMgr ripemgr.RipeMgr, geo geomgr.GeoMgr) (ProbeService, error) {
	return &ProbeServiceImpl{
		DBMgr:   dbMgr,
		RipeMgr: ripeMgr,
		GeoMgr:  geo,
	}, nil
}

func (srv *ProbeServiceImpl) RequestProbes() error {
	dbc := srv.DBMgr.GetDB()

	places, err := srv.findProbes(false)
	if err != nil {
		return err
	}
	placesAnchors, err := srv.findProbes(false)
	if err != nil {
		return err
	}

	locsList := []*models.Location{}
	dbc.Order(clause.OrderByColumn{Column: clause.Column{Name: "country"}}).
		Find(&locsList)
	for _, location := range locsList {
		log.WithFields(log.Fields{
			"country": location.Country,
			"iata":    location.IataCode,
		}).Info("Get probes for airport")

		nearestAnchorProbeIDs := measurements.FindNearest(placesAnchors,
			measurements.Place{Latitude: location.Latitude, Longitude: location.Longitude},
			config.NewConfig().GetInt("RIPE_ANCHOR_PROBES_PER_AIRPORT"))

		nearestProbeIDs := measurements.FindNearest(places,
			measurements.Place{Latitude: location.Latitude, Longitude: location.Longitude},
			config.NewConfig().GetInt("RIPE_PROBES_PER_AIRPORT"))

		nearestProbeIDs = append(nearestProbeIDs, nearestAnchorProbeIDs...)
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
		err = dbc.Model(location).Association("Probes").Replace(location.Probes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *ProbeServiceImpl) findProbes(isAnchor bool) ([]measurements.Place, error) {
	var places []measurements.Place
	err := srv.DBMgr.GetDB().Model(&models.Probe{}).
		Where(&models.Probe{IsAnchor: isAnchor, Status: "Connected"}).Find(&places).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get places from db")

		return nil, err
	}
	return places, nil
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
		}).Error("Update")

		return
	}

	log.Println("Probes successfully updated")
}

func (srv *ProbeServiceImpl) ImportProbes() {
	opts := make(map[string]string)
	opts["status_name"] = "Connected"
	opts["is_public"] = "true"
	opts["sort"] = "id"

	probes, err := srv.RipeMgr.GetProbes(opts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get all probes from ripe")
		return
	}

	probesDB := []*models.Probe{}
	for _, v := range probes { //nolint:gocritic
		newProbe := &models.Probe{
			ProbeID:     v.ID,
			CountryCode: v.CountryCode,
			Status:      v.Status.Name,
			IsAnchor:    v.IsAnchor,
			IsPublic:    v.IsPublic,
		}

		srv.fillCoordinates(v, newProbe)
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
			DoUpdates: clause.AssignmentColumns([]string{"latitude", "longitude", "status", "swap_coordinates", "is_anchor", "is_public"}),
		}).Create(probesDB).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to insert probes")
	}
}

func (srv *ProbeServiceImpl) fillCoordinates(v atlas.Probe, newProbe *models.Probe) { // nolint:gocritic
	if v.Geometry.Type == point {
		newProbe.Latitude = v.Geometry.Coordinates[0]
		newProbe.Longitude = v.Geometry.Coordinates[1]
		codeCountry := srv.GeoMgr.FindCountry(newProbe.Latitude, newProbe.Longitude)
		if codeCountry != v.CountryCode {
			codeCountry = srv.GeoMgr.FindCountry(newProbe.Longitude, newProbe.Latitude)
			if codeCountry == v.CountryCode {
				newProbe.Latitude = v.Geometry.Coordinates[1]
				newProbe.Longitude = v.Geometry.Coordinates[0]
				newProbe.SwapCoordinates = true
				log.WithFields(log.Fields{
					"CountryCode": v.CountryCode,
					"Latitude":    newProbe.Latitude,
					"Longitude":   newProbe.Longitude,
					"order:":      "wrong",
				}).Warn("coordinates was in wrong order")
			}
		}
	}
}
