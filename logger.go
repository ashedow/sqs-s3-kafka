package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

var LOG = log.NewJSONLogger(os.Stdout)

func logInfo(message interface{}, keyvals ...interface{}) {
	LOG.Log(
		append(keyvals,
			"level", "INFO",
			"message", message,
		)...,
	)
}

func logError(err error) {
	if err != nil {
		LOG.Log(
			"level", "ERROR",
			"message", err,
		)
	}
}
