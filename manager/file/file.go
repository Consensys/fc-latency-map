package file

import (
	"io/fs"
	"os"

	log "github.com/sirupsen/logrus"
)

func Create(fn string, fullJSON []byte) {
	if err := os.WriteFile(fn, fullJSON, fs.ModePerm); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create file")

		return
	}
}
