package ripemgr

import (
	"log"

	atlas "github.com/keltia/ripe-atlas"
)

type RipeMgrImpl struct {
	c         *atlas.Client
}

func NewRipeImpl(APIKey string) (RipeMgr, error) {
	cfgs := []atlas.Config{}
	cfgs = append(cfgs, atlas.Config{
		APIKey: APIKey,
	})
	c, err := atlas.NewClient(cfgs...)
	if err != nil {
		log.Println("Connecting to Ripe Atlas API", err)
		return nil, err
	}
	ver := atlas.GetVersion()
	log.Println("api version ", ver)
	
	return &RipeMgrImpl{
		c:  c,
	}, nil
}

func (fMgr *RipeMgrImpl) GetProbe(id int) (*atlas.Probe, error) {
	return fMgr.c.GetProbe(id)
}

func (fMgr *RipeMgrImpl) GetProbes(opts map[string]string) ([]atlas.Probe, error)  {
	probes, err := fMgr.c.GetProbes(opts)
	if err != nil {
		return nil, err
	}

	return probes, nil
}
