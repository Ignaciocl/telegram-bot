package utils

import log "github.com/sirupsen/logrus"

func LogError(err error, msg string) {
	if err != nil {
		log.Errorf("%s - %+v", msg, err)
	}
}
