package models

import (
	"bindolabs/sevenrooms/db"
	"bindolabs/sevenrooms/log"
	"time"

	"bindolabs/sevenrooms/gatewaymodels"

	"github.com/jinzhu/gorm"
)

type LineItem struct {
	ID                uint `gorm:"primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	GatewayLineItemID int     `gorm:"column:gateway_line_item_id; type:int(11) ;" json:"gateway_line_item_id"`
	TableID           uint    `gorm:"column:table_id; type:int(11) ;" json:"table_id"`
	Name              string  `gorm:"column:name; type:varchar(255) ;" json:"name"`
	Price             float64 `gorm:"column:price; type:decimal(16,2) ;default:'0.00'" json:"pirce"`
	NetTotal          float64 `gorm:"column:net_total; type:decimal(16,2) ;default:'0.00'" json:"net_total"`
	Qty               float64 `gorm:"column:qty; type:decimal(12,4) ;" json:"qty"`
	ListingBarcode    string  `gorm:"column:listing_barcode; type:varchar(255) ;" json:"listing_barcode"`
}

func (*LineItem) TableName() string {
	return "line_items"
}

func (l *LineItem) GetListingBarcode() {
	var listingSnapshot gatewaymodels.ListingSnapshot
	var err error
	if l.GatewayLineItemID == 0 {
		return
	}
	if err = db.GatewayDB.Model(&listingSnapshot).Where("line_item_id = ?", l.GatewayLineItemID).Find(&listingSnapshot).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			l.ListingBarcode = listingSnapshot.ListingBarcode.V()
			return
		}
		log.Logger.Errorf("find GatewayDB LineItems err: %s", err)
		return
	}
}
