package ripeapi

type RipeAPIImpl struct {
	addr string
	APIKey string
}

func NewRipeAPIImpl(addr string, APIKey string) RipeAPIImpl {
	return RipeAPIImpl{
		addr:   addr,
		APIKey:	APIKey,
	}
}

func (fMgr *RipeAPIImpl) GetProbes() error {
	return nil
}

func (fMgr *RipeAPIImpl) GetMeasurements() error {
	return nil
}