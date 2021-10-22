package probes

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/config"

	"github.com/ConsenSys/fc-latency-map/manager/geomgr"

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

const (
	point = "Point"
)

func NewProbeServiceImpl(dbMgr db.DatabaseMgr, ripeMgr ripemgr.RipeMgr, geo geomgr.GeoMgr) (ProbeService, error) {
	return &ProbeServiceImpl{
		DBMgr:   dbMgr,
		RipeMgr: ripeMgr,
		GeoMgr:  geo,
	}, nil
}

func (srv *ProbeServiceImpl) RequestProbes() error {
	dbc := srv.DBMgr.GetDB()

	places, err := srv.findProbesAsPlaces(&models.Probe{IsAnchor: false, Status: "Connected"})
	if err != nil {
		return err
	}
	placesAnchors, err := srv.findProbesAsPlaces(&models.Probe{IsAnchor: true, Status: models.StatusConnected})
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
		}).Info("get probes for location")

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

func (srv *ProbeServiceImpl) findProbesAsPlaces(query interface{}) ([]measurements.Place, error) {
	var places []measurements.Place
	err := srv.DBMgr.GetDB().
		Model(&models.Probe{}).
		Where(query).
		Find(&places).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get places from db")

		return nil, err
	}
	return places, nil
}
func (srv *ProbeServiceImpl) findProbes(query interface{}) ([]*models.Probe, error) {
	var ps []*models.Probe
	err := srv.DBMgr.GetDB().
		Model(&models.Probe{}).
		Where(query).
		Find(&ps).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get places from db")

		return nil, err
	}
	return ps, nil
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

func (srv *ProbeServiceImpl) Update() bool {
	err := srv.RequestProbes()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Update")

		return false
	}
	log.Println("Probes successfully updated")
	return true
}

func (srv *ProbeServiceImpl) ImportProbes() bool {
	opts := make(map[string]string)
	opts["status_name"] = "Connected"
	opts["is_public"] = "true"
	opts["sort"] = "id"

	probes, err := srv.RipeMgr.GetProbes(opts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("get all probes from ripe")
		return false
	}

	probesToSave := []*models.Probe{}
	for _, v := range probes { //nolint:gocritic
		newProbe := &models.Probe{
			ProbeID:     v.ID,
			CountryCode: v.CountryCode,
			Status:      models.Status(v.Status.Name),
			IsAnchor:    v.IsAnchor,
			IsPublic:    v.IsPublic,
			AddressV4:   v.AddressV4,
			AddressV6:   v.AddressV6,
		}

		if v.Geometry.Type == point {
			newProbe.RipeLatitude = v.Geometry.Coordinates[1]
			newProbe.RipeLongitude = v.Geometry.Coordinates[0]
		}
		probesToSave = append(probesToSave, newProbe)
	}

	dbc := srv.DBMgr.GetDB()

	// update all the rows
	err = dbc.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Model(&models.Probe{}).
		Updates(models.Probe{Status: models.StatusDisconnected}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to update deprecated probes")
	}

	srv.upsertProbes(dbc, probesToSave, []string{"status", "is_anchor", "is_public", "address_v4", "address_v6"})

	srv.upsertProbesCoordinates()

	return true
}

func (srv *ProbeServiceImpl) upsertProbes(dbc *gorm.DB, probesDB []*models.Probe, updtColumns []string) {
	err := dbc.Model(&models.Probe{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "probe_id"}},
			DoUpdates: clause.AssignmentColumns(updtColumns),
		}).Create(probesDB).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to insert probes")
	}
}

func (srv *ProbeServiceImpl) fixCoordinates(p *models.Probe) {
	ip := p.AddressV4
	if ip == "" {
		ip = p.AddressV6
	}

	p.CoordinatesStatus = models.CoordinatesStatusOk
	geo, err := srv.GeoMgr.IPGeolocation(ip)
	if err != nil {
		return
	}
	distance := measurements.Haversine(p.RipeLatitude, p.RipeLongitude, geo.Latitude, geo.Longitude)

	const distanceBtwCoodinates = 1000
	if distance > distanceBtwCoodinates {
		p.Longitude = geo.Longitude
		p.Latitude = geo.Latitude
		log.WithFields(log.Fields{
			"distance:":       distance,
			"probe":           p.ProbeID,
			"RipeCountryCode": p.CountryCode,
			"RipeLatitude":    p.RipeLatitude,
			"RipeLongitude":   p.RipeLongitude,
			"Country":         geo.Country,
			"Latitude":        geo.Latitude,
			"Longitude":       geo.Longitude,
		}).Warn("faraway places coordinates")
	} else {
		p.Latitude = p.RipeLatitude
		p.Longitude = p.RipeLongitude
	}
}

func (srv *ProbeServiceImpl) upsertProbesCoordinates() {
	p, _ := srv.findProbes(&models.Probe{
		Status:            models.StatusConnected,
		CoordinatesStatus: models.CoordinatesStatusUnknown,
	})
	for _, probe := range p {
		srv.fixCoordinates(probe)
	}
	srv.upsertProbes(srv.DBMgr.GetDB(), p, []string{"latitude", "longitude", "coordinates_status"})
}
