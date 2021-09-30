package measurements

import (
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/spatial/vptree"
	"gorm.io/gorm"
)

func FindNearest(q Place, amount int, table string, dbi *gorm.DB) []int {
	var places []Place
	err := dbi.Debug().Table(table).Where("deleted_at IS null").Find(&places).Error
	if err != nil {
		log.WithFields(log.Fields{
			"table": table,
			"error": err,
		}).Error("find latitude/longitude from db")
		return nil
	}
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
		log.Fatal(err)
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

		"Probe IDs": ids,
		"lat":       q.Latitude,
		"lon":       q.Longitude,
		"amount":    amount,
	}).Info("FindNearest locations")

	return ids
}
