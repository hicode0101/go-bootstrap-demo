package services

import (
	"go-bootstrap-demo/models"

	logger "github.com/hicode0101/go-logger"
)

var (
	ConfigS *models.Config
)

func init() {
	logger.Debug("services init")

	LoadConfig()
}
