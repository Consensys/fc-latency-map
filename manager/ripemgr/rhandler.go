package ripemgr

import (
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/models"

	atlas "github.com/keltia/ripe-atlas"
)

type Handler struct {
	Service *Service
}

func NewHandler() *Handler {
	conf := config.NewConfig()

	var apiKey = conf.GetString("RIPE_API_KEY")
	cfg := atlas.Config{
		APIKey: apiKey,
	}
	ripe, err := atlas.NewClient(cfg)
	if err != nil {
		log.Fatalf("connecting with ripe failed: %s", err)
	}

	mSer := NewServiceImpl(conf, ripe)

	return &Handler{
		Service: &mSer,
	}
}

func (h *Handler) CreateMeasurement(miners []*models.Miner, value string) ([]*atlas.Measurement, error) {
	return (*h.Service).createMeasurements(miners, value)

}

func (h *Handler) GetMeasurementResults(measures map[int]int) ([]atlas.MeasurementResult, error) {
	return (*h.Service).getMeasurementResults(measures)
}
