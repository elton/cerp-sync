package main

import (
	"runtime"

	"github.com/elton/cerp-sync/cron"
	"github.com/elton/cerp-sync/models"
	"github.com/elton/cerp-sync/utils/batch"
	"github.com/elton/cerp-sync/utils/logger"
)

func main() {
	runtime.GOMAXPROCS(1)

	var shop *models.Shop
	shops, err := shop.GetAllShops()
	if err != nil {
		logger.Error.Println(err)
	}

	if len(*shops) <= 0 {
		batch.InitializeData()
		cron.SyncData()
		select {}
	} else {
		cron.SyncData()
		select {}
	}
}
