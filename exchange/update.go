package exchange

import (
	"fmt"
	"log"
	"time"

	"autumn-2021-intern-assignment/utils/config"
)

func UpdateWithLoad(conf config.Exchanger, currPath string, skip, load bool) {
	if !load {
		log.Println("Currencies Update started")
		err := Update(conf)
		log.Println("Currencies Update done")
		if err != nil {
			log.Println(fmt.Errorf("updating currencies: %w", err))

			err = Load(currPath)
			if err != nil {
				log.Println(fmt.Errorf("load currencies: %w", err))
			} else {
				log.Println("Currencies were loaded from file")
			}
		} else {
			err = Store(currPath)
			if err != nil {
				log.Println(fmt.Errorf("store currencies: %w", err))
			}
		}
		if !skip {
			go func() {
				for {
					time.Sleep(time.Second * time.Duration(conf.Every))
					err = Update(conf)
					if err != nil {
						log.Println(fmt.Errorf("updating currencies: %w", err))
						err = Update(conf)
					}
				}
			}()
		}
	} else {
		err := Load(currPath)
		if err != nil {
			log.Println(fmt.Errorf("load currencies: %w", err))
		} else {
			log.Println("Currencies were loaded from file")
		}
	}
}
