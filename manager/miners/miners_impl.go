package miners

import (
	"log"

	"github.com/ConsenSys/fc-latency-map/manager/addresses"
	"github.com/ConsenSys/fc-latency-map/manager/db"
	fmgr "github.com/ConsenSys/fc-latency-map/manager/filecoinmgr"
	"github.com/ConsenSys/fc-latency-map/manager/models"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

type MinerServiceImpl struct {
	Conf  *viper.Viper
	DbMgr *db.DatabaseMgr
	FMgr  *fmgr.FilecoinMgr
}

func NewMinerServiceImpl(conf *viper.Viper, dbMgr *db.DatabaseMgr, fMgr *fmgr.FilecoinMgr) MinerService {
	return &MinerServiceImpl{
		Conf:  conf,
		DbMgr: dbMgr,
		FMgr:  fMgr,
	}
}

func (srv *MinerServiceImpl) GetAllMiners() []*models.Miner {
	var miners []*models.Miner
	(*srv.DbMgr).GetDb().Find(&miners)
	for _, miner := range miners {
		log.Printf("Miner address: %s - ip: %s\n", miner.Address, miner.Ip)
	}
	return miners
}

func (srv *MinerServiceImpl) ParseMiners(offset uint) []*models.Miner {
	blockHeight, err := (*srv.FMgr).GetBlockHeight()
	if err != nil {
		log.Fatalf("get block failed: %s", err)
		return []*models.Miner{}
	}
	log.Printf("blockHeight: %+v\n", blockHeight)
	deals, err := (*srv.FMgr).GetVerifiedDeals(blockHeight, offset)
	if err != nil {
		log.Fatalf("get block failed: %s", err)
		return []*models.Miner{}
	}
	return srv.parseMinersFromDeals(deals)
}

func (srv *MinerServiceImpl) parseMinersFromDeals(deals []fmgr.VerifiedDeal) []*models.Miner {
	var miners = []*models.Miner{}
	for _, deal := range deals {
		provider := deal.Provider
		address := provider.String()
		minerInfo, err := (*srv.FMgr).GetMinerInfo(provider)
		if err != nil {
			log.Printf("unable to get miner info: %s. skip...", address)
			continue
		}
		miners = append(miners, &models.Miner{
			Address: address,
			Ip:      getMinerIp(minerInfo),
		})
	}
	if len(miners) > 0 {
		srv.upsertMinersInDb(miners)
		for _, miner := range miners {
			log.Printf("Miner address: %s - ip: %s\n", miner.Address, miner.Ip)
		}
	} else {
		log.Printf("No miner parsed")
	}
	return miners
}

func getMinerIp(minerInfo miner.MinerInfo) string {
	ips := addresses.IpAddress(addresses.MultiAddrs(minerInfo.Multiaddrs))
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

func (srv *MinerServiceImpl) upsertMinersInDb(miners []*models.Miner) {
	(*srv.DbMgr).GetDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"ip"}),
	}).Create(&miners)
}

func (srv *MinerServiceImpl) ParseMinersByBlockHeight(height int64) []*models.Miner {
	deals, err := (*srv.FMgr).GetVerifiedDealsByBlockHeight(abi.ChainEpoch(height))
	if err != nil {
		log.Fatalf("get block failed: %s", err)
		return []*models.Miner{}
	}
	return srv.parseMinersFromDeals(deals)
}
