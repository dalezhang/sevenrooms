package restaurantmodels

import (
	"bindolabs/sevenrooms/db"
	"bindolabs/sevenrooms/log"
	"time"

	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
)

type Party struct {
	CreatedAt    time.Time       `gorm:"column:created_at; type:datetime ;" json:"created_at"`
	DeprecatedAt *time.Time      `gorm:"column:deprecated_at; type:datetime ;" json:"deprecated_at"`
	ID           optional.Int    `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	IsBatch      optional.Bool   `gorm:"column:is_batch; type:tinyint(1) ;" json:"is_batch"`
	IsParent     optional.Bool   `gorm:"column:is_parent; type:tinyint(1) ;" json:"is_parent"`
	Note         optional.String `gorm:"column:note; type:text ;" json:"note"`
	ParentUUID   []byte          `gorm:"column:parent_uuid; type:binary(16) ;" json:"parent_uuid"`
	PartyGroupID optional.Int    `gorm:"column:party_group_id; type:int(11) ;" json:"party_group_id"`
	People       optional.Int    `gorm:"column:people; type:int(11) ;" json:"people"`
	ReopenedAt   *time.Time      `gorm:"column:reopened_at; type:datetime ;" json:"reopened_at"`
	SeatedAt     *time.Time      `gorm:"column:seated_at; type:datetime ;" json:"seated_at"`
	TableID      optional.Int    `gorm:"column:table_id; type:int(11) ;" json:"table_id"`
	TableSplit   optional.Int    `gorm:"column:table_split; type:int(11) ;" json:"table_split"`
	UnseatedAt   *time.Time      `gorm:"column:unseated_at; type:datetime ;" json:"unseated_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at; type:datetime ;" json:"updated_at"`
	UUID         []byte          `gorm:"column:uuid; type:binary(16) ;" json:"uuid"`
}

// TableName sets the insert table name for this struct type
func (p *Party) TableName() string {
	return "parties"
}

func (p *Party) GetTable() (err error, tb Table) {
	if err = db.RestDB.Model(&tb).Find(&tb, p.TableID.V()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("Table err: %s", err)
	}
	return
}
