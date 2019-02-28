package gatewaymodels

import (
	"github.com/imiskolee/optional"
)

type ListingSnapshot struct {
	ID             optional.Int    `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	LineItemID     optional.Int    `gorm:"column:line_item_id; type:int(11) ;" json:"line_item_id"`
	ListingBarcode optional.String `gorm:"column:listing_barcode; type:varchar(255) ;" json:"listing_barcode"`
}

// TableName sets the insert table name for this struct type
func (l *ListingSnapshot) TableName() string {
	return "listing_snapshots"
}
