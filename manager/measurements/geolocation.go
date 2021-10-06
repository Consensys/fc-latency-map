package measurements

import (
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/spatial/vptree"
)

func FindNearest(places []Place, q Place, amount int) []int {
	if len(places) == 0 {
		return []int{}
	}
	var comparables []vptree.Comparable
	for _, place := range places {
		comparables = append(comparables, place)
	}
	const effort = 5
	t, err := vptree.New(comparables, effort, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"lat":    q.Latitude,
			"lon":    q.Longitude,
			"amount": amount,
		}).Error("findNearest locations vptree")

		return []int{}
	}

	// Find the five closest probes to the residence.
	keep := vptree.NewNKeeper(amount)
	t.NearestSet(keep, q)

	var ids []int
	for _, c := range keep.Heap {
		p := c.Comparable.(Place)
		ids = append(ids, p.ID)
	}

	log.WithFields(log.Fields{
		"IDs":    ids,
		"lat":    q.Latitude,
		"lon":    q.Longitude,
		"amount": amount,
	}).Info("FindNearest locations")

	return ids
}
