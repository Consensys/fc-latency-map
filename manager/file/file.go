package file

import (
	"io/fs"
	"os"

	log "github.com/sirupsen/logrus"
)

func Create(fn string, fullJson []byte) {
	err := os.WriteFile(fn, fullJson, fs.ModePerm)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Create file")
		return
	}
}
