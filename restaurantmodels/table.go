package restaurantmodels

import (
	"time"

	"github.com/imiskolee/optional"
)

type Table struct {
	AllowMultipleParties optional.Bool    `gorm:"column:allow_multiple_parties; type:tinyint(1) ;default:'1'" json:"allow_multiple_parties"`
	CreatedAt            time.Time        `gorm:"column:created_at; type:datetime ;" json:"created_at"`
	DefaultTurnTime      optional.Int     `gorm:"column:default_turn_time; type:int(11) ;" json:"default_turn_time"`
	DeprecatedAt         *time.Time       `gorm:"column:deprecated_at; type:datetime ;" json:"deprecated_at"`
	Height               optional.Int     `gorm:"column:height; type:int(11) ;" json:"height"`
	ID                   optional.Int     `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	Name                 optional.String  `gorm:"column:name; type:varchar(255) ;" json:"name"`
	PositionX            optional.Int     `gorm:"column:position_x; type:int(11) ;" json:"position_x"`
	PositionY            optional.Int     `gorm:"column:position_y; type:int(11) ;" json:"position_y"`
	RoomID               optional.Int     `gorm:"column:room_id; type:int(11) ;" json:"room_id"`
	Rotation             optional.Int     `gorm:"column:rotation; type:int(11) ;" json:"rotation"`
	Scale                optional.Float64 `gorm:"column:scale; type:float ;" json:"scale"`
	Seats                optional.Int     `gorm:"column:seats; type:int(11) ;" json:"seats"`
	ServerID             optional.Int     `gorm:"column:server_id; type:int(11) ;" json:"server_id"`
	Shape                optional.String  `gorm:"column:shape; type:varchar(255) ;" json:"shape"`
	Sides                optional.Int     `gorm:"column:sides; type:int(11) ;" json:"sides"`
	StoreID              optional.Int     `gorm:"column:store_id; type:int(11) ;" json:"store_id"`
	UpdatedAt            time.Time        `gorm:"column:updated_at; type:datetime ;" json:"updated_at"`
	UUID                 optional.String  `gorm:"column:uuid; type:varchar(191) ;" json:"uuid"`
	Width                optional.Int     `gorm:"column:width; type:int(11) ;" json:"width"`
}

// TableName sets the insert table name for this struct type
func (t *Table) TableName() string {
	return "tables"
}
