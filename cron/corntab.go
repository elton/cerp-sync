package cron

import (
	"strconv"
	"time"

	"github.com/elton/cerp-sync/broker"
	"github.com/elton/cerp-sync/config"
	"github.com/elton/cerp-sync/models"
	"github.com/elton/cerp-sync/utils/logger"
	"github.com/robfig/cron"
)

// SyncData synchron all the data of shop and order.
func SyncData() {
	logger.Info.Println("Begin sync task.")
	// Sync store information
	c := cron.New()

	c.AddFunc("@hourly", func() {
		_, err := getShops()
		if err != nil {
			logger.Error.Println(err)
			return
		}
	})

	// Sync order information.
	c.AddFunc("* */2 * * * ?", func() {
		var shop models.Shop

		shops, err := shop.GetAllShops()
		if err != nil {
			logger.Error.Println(err)
			return
		}

		for _, shop := range *shops {
			if err := getOrders(shop.Code); err != nil {
				logger.Error.Println(err)
				return
			}
		}
	})

	c.Start()

}

// getShops save all shop information.
func getShops() (*[]models.Shop, error) {
	var (
		shop               models.Shop
		shops, shopCreated *[]models.Shop
		lastUpdateAt       time.Time
		err                error
	)

	if lastUpdateAt, err = shop.GetLastUpdatedAt(); err != nil {
		return nil, err
	}

	if shops, err = broker.GetShops("1", "20", lastUpdateAt); err != nil {
		return nil, err
	}

	if len(*shops) > 0 {
		if shopCreated, err = shop.SaveAll(shops); err != nil {
			return nil, err
		}
		logger.Info.Printf("Save %d shops information\n", len(*shopCreated))
	}

	return shops, nil
}

// getOrders save all the orders of specified shop.
func getOrders(shopCode string) error {
	var (
		orderDb              models.Order
		orders, orderCreated *[]models.Order
		lastUpdateAt         time.Time
		totalOrder           int
		err                  error
	)

	pgSize, _ := strconv.Atoi(config.Config("PAGE_SIZE"))

	if lastUpdateAt, err = orderDb.GetLastUpdatedAt(shopCode); err != nil {
		return err
	}

	if totalOrder, err = broker.GetTotalOfOrders(shopCode, lastUpdateAt); err != nil {
		return err
	}

	totalPg := totalOrder / pgSize
	if totalOrder%pgSize != 0 {
		totalPg = totalPg + 1
	}

	logger.Info.Printf("Total Order: %d, page size: %d, total page: %d", totalOrder, pgSize, totalPg)

	for i := 0; i < totalPg; i++ {
		if orders, err = broker.GetOrders(strconv.Itoa(i+1), strconv.Itoa(pgSize), shopCode, lastUpdateAt); err != nil {
			return err
		}

		if len(*orders) > 0 {
			if orderCreated, err = orderDb.SaveAll(orders); err != nil {
				return err
			}
			logger.Info.Printf("Save (shop %s) %d orders information, page %d of %d\n", shopCode, len(*orderCreated), i+1, totalPg)
		}
	}
	return nil
}
