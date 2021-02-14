package batch

import (
	"os"
	"strconv"
	"time"

	"github.com/elton/cerp-sync/broker"
	"github.com/elton/cerp-sync/models"
	"github.com/elton/cerp-sync/utils/logger"
	"github.com/joho/godotenv"
)

// InitializeData save all shops and orders data in the database.
func InitializeData() {

	shops, err := getShops()
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
}

// getShops save all shop information.
func getShops() (*[]models.Shop, error) {
	var (
		shop               models.Shop
		shops, shopCreated *[]models.Shop
		// layout             string = "2006-01-02 15:04:05"
		begin time.Time
		err   error
	)

	// begin, err := time.Parse(layout, "0001-01-01 00:00:00")
	// if err != nil {
	// 	return nil, err
	// }

	if shops, err = broker.GetShops("1", "20", begin); err != nil {
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
		totalOrder           int
		// layout               string = "2006-01-02 15:04:05"
		begin time.Time
		err   error
	)
	godotenv.Load()
	pgSize, _ := strconv.Atoi(os.Getenv("PAGE_SIZE"))

	if totalOrder, err = broker.GetTotalOfOrders(shopCode, begin); err != nil {
		return err
	}

	totalPg := totalOrder / pgSize
	if totalOrder%pgSize != 0 {
		totalPg = totalPg + 1
	}

	logger.Info.Printf("Total Order: %d, page size: %d, total page: %d", totalOrder, pgSize, totalPg)

	for i := 0; i < totalPg; i++ {
		if orders, err = broker.GetOrders(strconv.Itoa(i+1), strconv.Itoa(pgSize), shopCode, begin); err != nil {
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
