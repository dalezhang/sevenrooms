package restaurantmodels

import (
	"bindolabs/optitable_middleware/db"
	"bindolabs/optitable_middleware/log"
	"time"

	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
)

type PartyGroup struct {
	ArrivedAt              *time.Time      `gorm:"column:arrived_at; type:datetime ;" json:"arrived_at"`
	BookingType            optional.Int    `gorm:"column:booking_type; type:int(11) ;" json:"booking_type"`
	ClearedAt              *time.Time      `gorm:"column:cleared_at; type:datetime ;" json:"cleared_at"`
	CreatedAt              time.Time       `gorm:"column:created_at; type:datetime ;" json:"created_at"`
	CustomerID             optional.Int    `gorm:"column:customer_id; type:int(11) ;" json:"customer_id"`
	DeprecatedAt           *time.Time      `gorm:"column:deprecated_at; type:datetime ;" json:"deprecated_at"`
	ID                     optional.Int    `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	Name                   optional.String `gorm:"column:name; type:varchar(255) ;" json:"name"`
	OrderID                optional.Int    `gorm:"column:order_id; type:int(11) ;" json:"order_id"`
	Overtime               optional.Int    `gorm:"column:overtime; type:int(11) ;default:'0'" json:"overtime"`
	PartySizeSegmentItemID optional.Int    `gorm:"column:party_size_segment_item_id; type:int(11) ;" json:"party_size_segment_item_id"`
	Phone                  optional.String `gorm:"column:phone; type:varchar(255) ;" json:"phone"`
	ReservationTime        *time.Time      `gorm:"column:reservation_time; type:datetime ;" json:"reservation_time"`
	SequenceNumber         optional.Int    `gorm:"column:sequence_number; type:int(11) ;" json:"sequence_number"`
	Status                 optional.String `gorm:"column:status; type:varchar(255) ;" json:"status"`
	StoreID                optional.Int    `gorm:"column:store_id; type:int(11) ;" json:"store_id"`
	TicketNumber           optional.String `gorm:"column:ticket_number; type:varchar(255) ;" json:"ticket_number"`
	TurnTime               optional.Int    `gorm:"column:turn_time; type:int(11) ;" json:"turn_time"`
	UpdatedAt              time.Time       `gorm:"column:updated_at; type:datetime ;" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (r *PartyGroup) TableName() string {
	return "party_groups"
}

func (r *PartyGroup) GetParties() (err error, pts []Party) {
	var pt Party
	if err = db.RestDB.Model(&pt).Where("`party_group_id` = ?", r.ID.V()).Find(&pts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("Party err: %s", err)
	}
	return
}
