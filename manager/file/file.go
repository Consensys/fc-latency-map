package file

import (
	"io/fs"
	"os"

	log "github.com/sirupsen/logrus"
)

func Create(fn string, fullJSON []byte) {
	err := os.WriteFile(fn, fullJSON, fs.ModePerm)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create file")
		return
	}
}
