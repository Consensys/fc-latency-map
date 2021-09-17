package ripemgr

import atlas "github.com/keltia/ripe-atlas"

type RipeMgr interface {

	// GetProbe return a probe by id
	GetProbe(id int) (*atlas.Probe, error)

	// GetProbes return probes list
	GetProbes(opts map[string]string) ([]atlas.Probe, error)

}
