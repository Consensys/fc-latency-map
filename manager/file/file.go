package file

import (
	"io/fs"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const layoutISO = "2006-01-02"
const jsonFolder = "data/export/"

// Create file
func Create(fn string, fullJSON []byte) {
	if err := os.WriteFile(jsonFolder+fn, fullJSON, fs.ModePerm); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create file")

		return
	}
}

// IsUpdated returns true when file exists and was created after the filename day
func IsUpdated(fn, date string) bool {
	file, err := os.Open(jsonFolder + fn)
	if err != nil {
		return false
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	t, err := time.Parse(layoutISO, date)
	if err != nil {
		return false
	}
	const oneDay = time.Hour * 24
	return stat.ModTime().After(t.Add(oneDay))
}
