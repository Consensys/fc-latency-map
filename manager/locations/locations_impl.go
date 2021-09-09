package locations

import (
	"github.com/ConsenSys/fc-latency-map/manager/db"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type LocationServiceImpl struct {
	DbMgr db.DatabaseMgr
}

func NewLocationServiceImpl(dbMgr db.DatabaseMgr) LocationService {
	return &LocationServiceImpl{
		DbMgr: dbMgr,
	}
}

func (srv *LocationServiceImpl) GetLocations() []*models.Location {
	var locations = []*models.Location{}
	srv.DbMgr.GetDb().Find(&locations)
	return locations
}

func (srv *LocationServiceImpl) GetLocation(location models.Location) models.Location {
	if err := srv.DbMgr.GetDb().Where(location).First(&location).Error; err != nil {
		return models.Location{}
	}
	return location
}

func (srv *LocationServiceImpl) AddLocation(location models.Location) models.Location {
	srv.DbMgr.GetDb().Create(&location) 
	return location
}

func (srv *LocationServiceImpl) DeleteLocation(location models.Location) bool {
	srv.DbMgr.GetDb().Delete(&location)
	return true
}
