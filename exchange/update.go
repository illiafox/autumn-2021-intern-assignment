package exchange

import (
	"time"

	"autumn-2021-intern-assignment/utils/config"
	"go.uber.org/zap"
)

func UpdateWithLoad(conf config.Exchanger, logger *zap.Logger) {

	Add("RUB", 1)

	if conf.Skip {
		return
	}

	if conf.Load {
		err := Load(conf.Path)
		if err != nil {
			logger.Error("loading currencies", zap.Error(err))
		} else {
			logger.Info("Done")
		}

		return
	}

	logger.Info("Currencies Update started")
	t := time.Now()
	err := Update(conf)

	if err != nil {
		logger.Error("updating currencies", zap.Duration("execution", time.Since(t)), zap.Error(err))

		logger.Info("loading from file")

		err = Load(conf.Path)
		if err != nil {
			logger.Error("loading currencies", zap.Error(err))
		} else {
			logger.Info("Done")
		}

	} else {
		logger.Info("Done", zap.Duration("execution", time.Since(t)))

		err = Store(conf.Path)
		if err != nil {
			logger.Error("writing currencies to file", zap.Error(err))
		}
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(conf.Every))
			err = Update(conf)
			if err != nil {
				logger.Error("updating currencies", zap.Error(err))
			}
		}
	}()

}
