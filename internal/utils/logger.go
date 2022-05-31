package utils

import (
	log "github.com/sirupsen/logrus"
)

func LogInfo(fields interface{}) {
	log.WithFields(
		log.Fields{
			// add fields if nessory
		}).Info(fields)
}
