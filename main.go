package main

import (
	_ "github.com/dimiro1/banner/autoload"

	"github.com/elton/cerp-sync/cron"
	"github.com/elton/cerp-sync/models"
	"github.com/elton/cerp-sync/utils/batch"
	"github.com/elton/cerp-sync/utils/logger"
)

func main() {
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
