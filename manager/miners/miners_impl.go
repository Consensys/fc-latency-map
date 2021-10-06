package miners

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"

	"github.com/ConsenSys/fc-latency-map/manager/addresses"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/geomgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
)

type MinerServiceImpl struct {
	Conf  *viper.Viper
	DBMgr db.DatabaseMgr
	FMgr  fmgr.FilecoinMgr
	GMgr  geomgr.GeoMgr
}

func NewMinerServiceImpl(conf *viper.Viper, dbMgr db.DatabaseMgr, fMgr fmgr.FilecoinMgr, gmgr geomgr.GeoMgr) MinerService {
	return &MinerServiceImpl{
		Conf:  conf,
		DBMgr: dbMgr,
		FMgr:  fMgr,
		GMgr:  gmgr,
	}
}

func (srv *MinerServiceImpl) GetAllMiners() []*models.Miner {
	var miners []*models.Miner
	srv.DBMgr.GetDB().Find(&miners)
	for _, m := range miners {
		log.Printf("Miner address: %s - ip: %s\n", m.Address, m.IP)
	}

	return miners
}

func (srv *MinerServiceImpl) ParseMinersByBlockOffset(offset int) []*models.Miner {
	blockHeight, err := (srv.FMgr).GetBlockHeight()
	if err != nil {
		log.Printf("GetBlockHeight failed: %s", err)
		return []*models.Miner{}
	}
	log.Printf("blockHeight: %+v\n", blockHeight)
	deals, err := (srv.FMgr).GetVerifiedDealsByBlockRange(blockHeight, offset)
	if err != nil {
		log.Printf("get Verified Deals By Block Range failed: %s", err)
		return []*models.Miner{}
	}

	return srv.parseMinersFromDeals(deals)
}

func (srv *MinerServiceImpl) parseMinersFromDeals(deals []fmgr.VerifiedDeal) []*models.Miner {
	miners := []*models.Miner{}
	for _, deal := range deals {
		provider := deal.Provider
		address := provider.String()
		minerInfo, err := (srv.FMgr).GetMinerInfo(provider)
		if err != nil {
			log.Printf("unable to get miner info: %s. skip...", address)

			continue
		}
		ip := getMinerIP(&minerInfo)
		lat, long := srv.getGeoLocation(ip)
		miners = append(miners, &models.Miner{
			Address:   address,
			IP:        ip,
			Latitude:  lat,
			Longitude: long,
		})
	}
	if len(miners) > 0 {
		srv.upsertMinersInDB(miners)
		for _, m := range miners {
			log.Printf("Miner address: %s - ip: %s\n", m.Address, m.IP)
		}
	} else {
		log.Printf("No miner parsed")
	}

	return miners
}

func (srv *MinerServiceImpl) getGeoLocation(ip string) (lat, long float64) {
	if ip != "" {
		split := strings.Split(ip, ",")

		return srv.GMgr.IPGeolocation(split[0])
	}

	return 0, 0
}

func getMinerIP(minerInfo *miner.MinerInfo) string {
	log.Printf("minerInfo.Multiaddrs: %s", minerInfo.Multiaddrs)
	ips := addresses.IPAddress(addresses.MultiAddrs(minerInfo.Multiaddrs))
	return strings.Join(ips, ",")
}

func (srv *MinerServiceImpl) upsertMinersInDB(miners []*models.Miner) {
	err := srv.DBMgr.GetDB().Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"ip", "latitude", "longitude"}),
	}).Create(&miners).Error
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("upsertMinersInDB")
	}
}

func (srv *MinerServiceImpl) ParseMinersByBlockHeight(height int64) []*models.Miner {
	deals, err := (srv.FMgr).GetVerifiedDealsByBlockHeight(abi.ChainEpoch(height))
	if err != nil {
		log.Printf("get Verified Deals By Block Height failed: %s", err)
		return []*models.Miner{}
	}
	return srv.parseMinersFromDeals(deals)
}

func (srv *MinerServiceImpl) ParseMinersByStateMarket() []*models.Miner {
	deals, err := (srv.FMgr).GetVerifiedDealsByStateMarket()
	if err != nil {
		log.Printf("get Verified Deals By Block Height failed: %s", err)
		return []*models.Miner{}
	}
	return srv.parseMinersFromDeals(deals)
}
