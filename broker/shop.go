package broker

import (
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/elton/cerp-sync/models"
	"github.com/elton/cerp-sync/utils/logger"
	"github.com/joho/godotenv"
)

// A Response struct to map the Entity Response
type Response struct {
	Success   bool   `json:"success"`
	ErrorDesc string `json:"errorDesc"`
	Total     int    `json:"total"`
	Shops     []Shop `json:"shops"`
}

// A Shop struct to map every shop information.
type Shop struct {
	ID         int    `json:"id"`
	Nick       string `json:"nick"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	CreateDate string `json:"create_date"`
	ModifyDate string `json:"modify_date"`
	Note       string `json:"note"`
	TypeName   string `json:"type_name"`
}

// GetShops returns the list of shops.
func GetShops(pgNum string, pgSize string, startDate time.Time) (*[]models.Shop, error) {
	var (
		shops          []models.Shop
		shop           models.Shop
		layout         string = "2006-01-02 15:04:05"
		responseObject Response
	)

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	request := make(map[string]interface{})
	request["appkey"] = os.Getenv("appKey")
	request["sessionkey"] = os.Getenv("sessionKey")
	request["method"] = "gy.erp.shop.get"
	request["page_no"] = pgNum
	request["page_size"] = pgSize
	if !startDate.IsZero() {
		request["modify_start_date"] = startDate.Format("2006-01-02 15:04:05")
	}

	if err := query(request, &responseObject); err != nil {
		return nil, err
	}

	logger.Info.Printf("Get %d shop information. \n", responseObject.Total)

	for _, _shop := range responseObject.Shops {
		shop.ID = node.Generate().Int64()
		shop.ShopID = _shop.ID
		shop.Name = _shop.Name
		shop.Nick = _shop.Nick
		shop.Code = _shop.Code
		shop.Note = _shop.Note
		shop.TypeName = _shop.TypeName

		if _shop.CreateDate != "" && _shop.CreateDate != "0000-00-00 00:00:00" {
			if shop.CreateDate, err = time.ParseInLocation(layout, _shop.CreateDate, time.Local); err != nil {
				return nil, err
			}
		}

		if _shop.ModifyDate != "" && _shop.ModifyDate != "0000-00-00 00:00:00" {
			if shop.ModifyDate, err = time.ParseInLocation(layout, _shop.ModifyDate, time.Local); err != nil {
				return nil, err
			}
		}

		shops = append(shops, shop)
	}

	return &shops, nil
}
