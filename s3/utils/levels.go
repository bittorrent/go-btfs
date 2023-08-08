package utils

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

func SetupLogLevels() {
	if _, set := os.LookupEnv("GOLOG_LOG_LEVEL"); !set {
		_ = logging.SetLogLevel("*", "INFO")

	} else {
		_ = logging.SetLogLevel("*", os.Getenv("GOLOG_LOG_LEVEL"))
	}
}
