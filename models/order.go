package models

import (
	"time"

	"github.com/elton/cerp-sync/utils/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// A Order struct to map the Entity Order
type Order struct {
	ID                   int64          `json:"id" gorm:"unique"`
	Code                 string         `json:"code" gorm:"primaryKey;size:256"`
	PlatformCode         string         `json:"platform_code" gorm:"size:256;index"`
	OrderTypeName        string         `json:"order_type_name" gorm:"size:256"`
	ShopName             string         `json:"shop_name" gorm:"size:256"`
	ShopCode             string         `json:"shop_code" gorm:"size:256;index"`
	VIPName              string         `json:"vip_name" gorm:"size:256;column:vip_name"`
	VIPCode              string         `json:"vip_code" gorm:"size:256;column:vip_code;index"`
	VIPRealName          string         `json:"vip_real_name" gorm:"size:256;column:vip_real_name"`
	AccountStatus        string         `json:"account_status" gorm:"size:256;index"`
	AccountAmount        float64        `json:"account_amount"`
	BusinessMan          string         `json:"business_man" gorm:"size:256"`
	Qty                  int8           `json:"qty"`
	Amount               float64        `json:"amount"`
	Payment              float64        `json:"payment"`
	WarehouseName        string         `json:"warehouse_name" gorm:"size:256"`
	WarehouseCode        string         `json:"warehouse_code" gorm:"size:256"`
	DeliveryState        int8           `json:"delivery_state" gorm:"index"`
	ExpressName          string         `json:"express_name" gorm:"size:256"`
	ExpressCode          string         `json:"express_code" gorm:"size:256;index"`
	ReceiverName         string         `json:"receiver_name" gorm:"size:256;index"`
	ReceiverMobile       string         `json:"receiver_mobile" gorm:"size:256;index"`
	ReceiverArea         string         `json:"receiver_area" gorm:"size:256"`
	ReceiverAddress      string         `json:"receiver_address" gorm:"size:512"`
	PlatformTradingState string         `json:"platform_trading_state" gorm:"size:256"`
	Deliverys            []Delivery     `json:"deliverys" gorm:"foreignKey:OrderCode;references:Code"`
	Details              []Detail       `json:"details" gorm:"foreignKey:OrderCode;references:Code"`
	Payments             []Payment      `json:"payments" gorm:"foreignKey:OrderCode;references:Code"`
	DealTime             time.Time      `json:"deal_time"`
	CreateTime           time.Time      `json:"create_time"`
	ModifyTime           time.Time      `json:"modify_date"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt            gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Delivery struct to map the Entity of the Delivery.
type Delivery struct {
	ID            int64          `json:"id" gorm:"primaryKey"`
	Delivery      bool           `json:"delivery" gorm:"comment:发货状态"`
	Code          string         `json:"code" gorm:"size:256;comment:发货单据号"`
	WarehouseName string         `json:"warehouse_name" gorm:"size:256"`
	WarehouseCode string         `json:"warehouse_code" gorm:"size:256;index"`
	ExpressName   string         `json:"express_name" gorm:"size:256"`
	ExpressCode   string         `json:"express_code" gorm:"size:256"`
	MailNo        string         `json:"mail_no" gorm:"size:256"`
	OrderCode     string         `json:"order_code" gorm:"size:256;index"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Detail struct to map the Entity of the item details.
type Detail struct {
	ID               int64          `json:"id" gorm:"primaryKey"`
	OID              string         `json:"oid" gorm:"size:256;comment:子订单号"`
	Qty              float64        `json:"qty"`
	Price            float64        `json:"price" gorm:"comment:实际单价"`
	Amount           float64        `json:"amount" gorm:"comment:实际金额"`
	Refund           int            `json:"refund" gorm:"comment:退款状态,0:未退款,1:退款成功,2:退款中"`
	Note             string         `json:"note"`
	PlatformItemName string         `json:"platform_item_name" gorm:"size:256;comment:平台规格名称"`
	PlatformSkuName  string         `json:"platform_sku_name" gorm:"size:256;comment:平台规格代码"`
	ItemCode         string         `json:"item_code" gorm:"size:256;index;comment:商品代码"`
	ItemName         string         `json:"item_name" gorm:"size:256;comment:商品名称"`
	ItemSimpleName   string         `json:"item_simple_name" gorm:"size:256;comment:商品简称"`
	PostFee          float64        `json:"post_fee" gorm:"comment:物流费用"`
	DiscountFee      float64        `json:"discount_fee" gorm:"comment:让利金额"`
	AmountAfter      float64        `json:"amount_after" gorm:"comment:让利后金额"`
	OrderCode        string         `json:"order_code" gorm:"size:256;index"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Payment struct to map the Entity of the payment.
type Payment struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	Payment     float64        `json:"payment" gorm:"comment:支付金额"`
	PayCode     string         `json:"pay_code" gorm:"size:256;comment:支付交易号"`
	PayTypeName string         `json:"pay_type_name" gorm:"size:256;comment:支付方式名称"`
	PayTime     time.Time      `json:"pay_time" gorm:"comment:支付时间"`
	OrderCode   string         `json:"order_code" gorm:"size:256;index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// SaveAll stores all specified the orders in the database.
func (o *Order) SaveAll(orders *[]Order) (*[]Order, error) {
	// 在冲突时，更新除主键以外的所有列到新值。
	if err := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetLastUpdatedAt get the last updated timestamp of the order.
func (o *Order) GetLastUpdatedAt(shopCode string) (time.Time, error) {
	var lastUpdateAt time.Time
	var layout string = "2006-01-02 15:04:05"
	if err := DB.Raw("SELECT orders.updated_at FROM orders WHERE orders.shop_code=? ORDER BY orders.updated_at DESC LIMIT 1", shopCode).Scan(&lastUpdateAt).Error; err != nil {
		rtime, err := time.Parse(layout, "0000-00-00 00:00:00")
		return rtime, err
	}
	logger.Info.Printf("Shop (%s): Order Last Updated: %v\n", shopCode, lastUpdateAt)
	return lastUpdateAt, nil
}
