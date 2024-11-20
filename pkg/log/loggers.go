package log

import (
	"smartDriver/internal/config"

	"go.uber.org/zap"
)

var (
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
)

func MustInit(cfg *config.Config) {
	var err error
	switch cfg.Server.Environment {
	case "development":
		Logger, err = zap.NewDevelopment()
	case "production":
		Logger, err = zap.NewProduction()
	}
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}

	SugaredLogger = Logger.Sugar()
	defer SugaredLogger.Sync()
}
