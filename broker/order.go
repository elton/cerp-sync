package broker

import (
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/elton/cerp-sync/models"
	"github.com/elton/cerp-sync/utils/logger"
	"github.com/joho/godotenv"
)

// A OResponse struct to map the Entity Response
type OResponse struct {
	Success    bool       `json:"success"`
	ErrorDesc  string     `json:"errorDesc"`
	Total      int        `json:"total"`
	Orders     []Order    `json:"orders"`
	Deliveries []Delivery `json:"deliverys"`
}

// A Order struct to map the Entity Order
type Order struct {
	Code                 string     `json:"code"`
	PlatformCode         string     `json:"platform_code"`
	OrderTypeName        string     `json:"order_type_name"`
	ShopName             string     `json:"shop_name"`
	ShopCode             string     `json:"shop_code"`
	VIPName              string     `json:"vip_name"`
	VIPCode              string     `json:"vip_code"`
	VIPRealName          string     `json:"vipRealName"`
	AccountStatus        string     `json:"accountStatus"`
	AccountAmount        float64    `json:"accountAmount"`
	BusinessMan          string     `json:"business_man"`
	Qty                  int8       `json:"qty"`
	Amount               float64    `json:"amount"`
	Payment              float64    `json:"payment"`
	WarehouseName        string     `json:"warehouse_name"`
	WarehouseCode        string     `json:"warehouse_code"`
	DeliveryState        int8       `json:"delivery_state"`
	ExpressName          string     `json:"express_name"`
	ExpressCode          string     `json:"express_code"`
	ReceiverName         string     `json:"receiver_phone"`
	ReceiverMobile       string     `json:"receiver_mobile"`
	ReceiverArea         string     `json:"receiver_area"`
	ReceiverAddress      string     `json:"receiver_address"`
	PlatformTradingState string     `json:"platform_trading_state"`
	Deliverys            []Delivery `json:"deliverys"`
	Details              []Detail   `json:"details"`
	Payments             []Payment  `json:"payments"`
	DealTime             string     `json:"dealtime"`
	CreateTime           string     `json:"createtime"`
	ModifyTime           string     `json:"modifytime"`
}

// Delivery struct to map the Entity of the Delivery.
type Delivery struct {
	Delivery      bool   `json:"delivery"`
	Code          string `json:"code"`
	WarehouseName string `json:"warehouse_name"`
	WarehouseCode string `json:"warehouse_code"`
	ExpressName   string `json:"express_name"`
	ExpressCode   string `json:"express_code"`
	MailNo        string `json:"mail_no"`
}

// Detail struct to map the Entity of the item details.
type Detail struct {
	OID              string  `json:"oid"`
	Qty              float64 `json:"qty"`
	Price            float64 `json:"price"`
	Amount           float64 `json:"amount"`
	Refund           int     `json:"refund"`
	Note             string  `json:"note"`
	PlatformItemName string  `json:"platform_item_name"`
	PlatformSkuName  string  `json:"platform_sku_name"`
	ItemCode         string  `json:"item_code"`
	ItemName         string  `json:"item_name"`
	ItemSimpleName   string  `json:"item_simple_name"`
	PostFee          float64 `json:"post_fee"`
	DiscountFee      float64 `json:"discount_fee"`
	AmountAfter      float64 `json:"amount_after"`
}

// Payment struct to map the Entity of the payment.
type Payment struct {
	Payment     float64 `json:"payment"`
	PayCode     string  `json:"payCode"`
	PayTypeName string  `json:"pay_type_name"`
	PayTime     string  `json:"payTime"`
}

// GetTotalOfOrders returns the total number of all orders.
func GetTotalOfOrders(shopCode string, startDate time.Time) (int, error) {
	err := godotenv.Load()
	if err != nil {
		return 0, err
	}

	request := make(map[string]interface{})

	request["appkey"] = os.Getenv("appKey")
	request["sessionkey"] = os.Getenv("sessionKey")
	request["method"] = "gy.erp.trade.get"
	request["shop_code"] = shopCode
	if !startDate.IsZero() {
		request["start_date"] = startDate.Format("2006-01-02 15:04:05")
	}
	request["date_type"] = 3

	var responseObject OResponse
	if err := query(request, &responseObject); err != nil {
		return 0, err
	}

	logger.Info.Printf("Get %d order information. \n", responseObject.Total)

	return responseObject.Total, nil
}

// GetOrders returns a list of all orders form specified shop.
func GetOrders(pgNum string, pgSize string, shopCode string, startDate time.Time) (*[]models.Order, error) {
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
	request["method"] = "gy.erp.trade.get"
	request["page_no"] = pgNum
	request["page_size"] = pgSize
	request["shop_code"] = shopCode
	if !startDate.IsZero() {
		request["start_date"] = startDate.Format("2006-01-02 15:04:05")
	}
	request["date_type"] = 3

	var (
		responseObject OResponse
		orders         []models.Order
		layout         string = "2006-01-02 15:04:05"
	)

	if err := query(request, &responseObject); err != nil {
		return nil, err
	}

	logger.Info.Printf("Get %d order information. \n", responseObject.Total)

	for _, _order := range responseObject.Orders {
		order := models.Order{}
		order.ID = node.Generate().Int64()
		order.Code = _order.Code
		order.PlatformCode = _order.PlatformCode
		order.OrderTypeName = _order.OrderTypeName
		order.ShopName = _order.ShopName
		order.ShopCode = _order.ShopCode
		order.VIPName = _order.VIPName
		order.VIPCode = _order.VIPCode
		order.VIPRealName = _order.VIPRealName
		order.AccountStatus = _order.AccountStatus
		order.AccountAmount = _order.AccountAmount
		order.BusinessMan = _order.BusinessMan
		order.Qty = _order.Qty
		order.Amount = _order.Amount
		order.Payment = _order.Payment
		order.WarehouseName = _order.WarehouseName
		order.WarehouseCode = _order.WarehouseCode
		order.DeliveryState = _order.DeliveryState
		order.ExpressName = _order.ExpressName
		order.ExpressCode = _order.ExpressCode
		order.ReceiverName = _order.ReceiverName
		order.ReceiverMobile = _order.ReceiverMobile
		order.ReceiverArea = _order.ReceiverArea
		order.ReceiverAddress = _order.ReceiverAddress
		order.PlatformTradingState = _order.PlatformTradingState

		for _, _delivery := range _order.Deliverys {
			delivery := models.Delivery{}
			delivery.ID = node.Generate().Int64()
			delivery.Delivery = _delivery.Delivery
			delivery.Code = _delivery.Code
			delivery.WarehouseName = _delivery.WarehouseName
			delivery.WarehouseCode = _delivery.WarehouseCode
			delivery.ExpressName = _delivery.ExpressName
			delivery.MailNo = _delivery.MailNo

			order.Deliverys = append(order.Deliverys, delivery)
		}

		for _, _detail := range _order.Details {
			detail := models.Detail{}
			detail.ID = node.Generate().Int64()
			detail.OID = _detail.OID
			detail.Qty = _detail.Qty
			detail.Price = _detail.Price
			detail.Amount = _detail.Amount
			detail.Refund = _detail.Refund
			detail.Note = _detail.Note
			detail.PlatformItemName = _detail.PlatformItemName
			detail.PlatformSkuName = _detail.PlatformSkuName
			detail.ItemCode = _detail.ItemCode
			detail.ItemName = _detail.ItemName
			detail.ItemSimpleName = _detail.ItemSimpleName
			detail.PostFee = _detail.PostFee
			detail.DiscountFee = _detail.DiscountFee
			detail.AmountAfter = _detail.AmountAfter

			order.Details = append(order.Details, detail)
		}

		for _, _payment := range _order.Payments {
			payment := models.Payment{}
			payment.ID = node.Generate().Int64()
			payment.Payment = _payment.Payment
			payment.PayCode = _payment.PayCode
			payment.PayTypeName = _payment.PayTypeName

			// 避免补发货订单没有支付信息，导致MySQL 报`Incorrect date value: '0000-00-00' for column`错误，
			// 参考 https://lefred.be/content/mysql-8-0-and-wrong-dates/ 设置MySQL参数

			if _payment.PayTime != "" && _payment.PayTime != "0000-00-00 00:00:00" {
				if payment.PayTime, err = time.ParseInLocation(layout, _payment.PayTime, time.Local); err != nil {
					return nil, err
				}
			}

			order.Payments = append(order.Payments, payment)
		}

		if _order.CreateTime != "" && _order.CreateTime != "0000-00-00 00:00:00" {
			if order.CreateTime, err = time.ParseInLocation(layout, _order.CreateTime, time.Local); err != nil {
				return nil, err
			}
		}
		if _order.ModifyTime != "" && _order.ModifyTime != "0000-00-00 00:00:00" {
			if order.ModifyTime, err = time.ParseInLocation(layout, _order.ModifyTime, time.Local); err != nil {
				return nil, err
			}
		}
		if _order.DealTime != "" && _order.DealTime != "0000-00-00 00:00:00" {
			if order.DealTime, err = time.ParseInLocation(layout, _order.DealTime, time.Local); err != nil {
				return nil, err
			}
		}

		orders = append(orders, order)
	}

	return &orders, nil
}
