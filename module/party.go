package module

import (
	"bindolabs/optitable_middleware/config"
	"bindolabs/optitable_middleware/db"
	"bindolabs/optitable_middleware/gatewaymodels"
	"bindolabs/optitable_middleware/log"
	"bindolabs/optitable_middleware/models"
	"bindolabs/optitable_middleware/restaurantmodels"
	"time"

	"fmt"

	"github.com/jinzhu/gorm"
)

type PartyMqMessage struct {
	StoreID     int    `json:"sid"`
	TableName   string `json:"table_name"`
	TableID     int    `json:"table_id"`
	PartyID     int    `json:"party_id"`
	People      int    `json:"people"`
	SeatedAt    int64  `json:"seated_at"`
	UnseatedAt  int64  `json:"unseated_at"`
	OrderNumber string `json:"order_number"`
}

func (p *PartyMqMessage) FilterPartyBySettingStores() (err error) {
	var isInSettingStore bool
	for _, store := range config.Conf.Setting.Stores {
		if p.StoreID == store.StoreID {
			isInSettingStore = true
		}
	}
	if isInSettingStore {
		return nil
	}
	return fmt.Errorf("store id: %d not in setting", p.StoreID)
}

func (p *PartyMqMessage) AnalyzeRecordFromMQ() (err error, table models.Table) {
	if err = db.DB.Model(table).Where("restaurant_party_id = ?", p.PartyID).Limit(1).Find(&table).Error; err == nil {
		if table.HasSync == true && table.Status == models.TableStatusClose {
			return
		}
		if p.UnseatedAt == 0 && p.SeatedAt != 0 {
			// 更新
			table.CheckOpenTime = p.SeatedAt
			table.CheckCloseTime = p.UnseatedAt
			fmt.Printf("\n +++CheckCloseTime: %d, %s", table.CheckCloseTime, table.Table)
			table.Status = models.TableStatusUpdate
			t := time.Now()
			table.RestaurantUpdatedAt = &t
			table.HasSync = false
			err = db.DB.Model(&table).Save(&table).Error
			if err != nil {
				log.Logger.Errorf("\n AnalyzeRecord save table err: %s", err)
				return
			}
			table.UpdateCheck()
		} else if p.UnseatedAt != 0 {
			// close check
			table.CheckOpenTime = p.SeatedAt
			table.CheckCloseTime = p.UnseatedAt
			fmt.Printf("\n +++CheckCloseTime: %d, %s", table.CheckCloseTime, table.Table)
			table.StoreID = p.StoreID
			table.Status = models.TableStatusClose
			t := time.Now()
			table.RestaurantUpdatedAt = &t
			table.HasSync = false
			err = db.DB.Model(&table).Save(&table).Error
			if err != nil {
				log.Logger.Errorf("\n AnalyzeRecord save table err: %s", err)
				return
			}
			// table.CloseCheck()
		}
	} else if err == gorm.ErrRecordNotFound {
		table.Table = p.TableName
		table.RestaurantTableID = p.TableID
		table.RestaurantPartyID = p.PartyID
		t := time.Now()
		table.RestaurantUpdatedAt = &t
		table.GuestCount = p.People
		table.CheckOpenTime = p.SeatedAt
		table.CheckCloseTime = p.UnseatedAt
		table.CheckRef = p.OrderNumber
		table.StoreID = p.StoreID
		table.Status = models.TableStatusCreate
		table.HasSync = false
		err = db.DB.Model(&table).Create(&table).Error
		if err != nil {
			log.Logger.Errorf("PartyID %d Create table err: %s", p.PartyID, err)
		}
		table.CreateCheck()
	} else {
		log.Logger.Errorf("CreateRecord err: %s", err)
		return err, table
	}
	return
}

func AnalyzeRecord(p *restaurantmodels.Party, pg *restaurantmodels.PartyGroup) (err error, table models.Table) {
	if pg.Status.V() == "check_dropped" {
		return
	}
	if err = db.DB.Model(table).Where("restaurant_party_id = ?", p.ID.V()).Limit(1).Find(&table).Error; err == nil {
		fmt.Printf("\n !!!! p.UnseatedAt = %v, p.SeatedAt = %v, table: %+v", p.UnseatedAt, p.SeatedAt, table)
		if table.HasSync == false {
			return
		}
		if table.HasSync == true && table.Status == models.TableStatusClose {
			return
		}
		if table.HasSync == true && table.Status == models.TableStatusVoid {
			return
		}
		if p.UnseatedAt == nil && p.SeatedAt != nil {
			// 换桌
			if p.TableID.V() != table.RestaurantTableID {
				var rTable restaurantmodels.Table
				if err = db.RestDB.Model(&rTable).Find(&rTable, p.TableID.V()).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
					}
					log.Logger.Errorf("CreateTable find table err: %s", err)
					return
				}
				table.OrgTable = table.Table
				table.Table = rTable.Name.V()
				table.RestaurantTableID = p.TableID.V()
				table.RestaurantUpdatedAt = &p.UpdatedAt
				table.GuestCount = p.People.V()
				if p.SeatedAt != nil {
					table.CheckOpenTime = p.SeatedAt.Unix()
				}
				if p.UnseatedAt != nil {
					table.CheckCloseTime = p.UnseatedAt.Unix()
				}
				table.Status = models.TableStatusMove
				table.HasSync = false
				err = db.DB.Model(&table).Save(&table).Error
				if err != nil {
					log.Logger.Errorf("PartyID %d AnalyzeRecord save table err: %s", p.ID.V(), err)
				}
				table.MoveTable()
				return
			}
			// 更新
			if p.SeatedAt.Unix() != table.CheckOpenTime || p.People.V() != table.GuestCount {
				table.CheckOpenTime = p.SeatedAt.Unix()
				table.Status = models.TableStatusUpdate
				table.RestaurantUpdatedAt = &p.UpdatedAt
				table.GuestCount = p.People.V()
				table.HasSync = false
				err = db.DB.Model(&table).Save(&table).Error
				if err != nil {
					log.Logger.Errorf("PartyID %d AnalyzeRecord save table err: %s", p.ID.V(), err)
					return
				}
				table.UpdateCheck()
			}
		} else if p.UnseatedAt != nil {
			// void check
			if p.UnseatedAt.Unix() != 0 && p.UnseatedAt.Unix() != table.CheckCloseTime {
				err, trans := table.GetAndSaveTransactions()
				if err != nil {
					return err, table
				}
				for _, t := range trans {
					t.VoidPaymentOrCloseCheck()
				}
				if err = db.DB.Model(table).Find(&table, table.ID).Error; err == nil {
					if table.HasSync == false {
						return err, table
					}
					if table.HasSync == true && table.Status == models.TableStatusClose {
						return err, table
					}
					if table.HasSync == true && table.Status == models.TableStatusVoid {
						return err, table
					}
				}

				var gOrder gatewaymodels.Order
				if err = db.GatewayDB.Model(&gOrder).Find(&gOrder, pg.OrderID.V()).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
					}
					log.Logger.Errorf("CreateTable find order err: %s", err)
					return err, table
				}
				if gOrder.State.V() != "voided" {
					return err, table
				}
				if p.SeatedAt != nil {
					table.CheckOpenTime = p.SeatedAt.Unix()
				}
				table.StoreID = pg.StoreID.V()
				table.Status = models.TableStatusVoid
				table.RestaurantUpdatedAt = &p.UpdatedAt
				table.HasSync = false
				err = db.DB.Model(&table).Save(&table).Error
				if err != nil {
					log.Logger.Errorf("\n AnalyzeRecord save table err: %s", err)
					return err, table
				}
				// table.VoidCheck()
			}
		}
	} else if err == gorm.ErrRecordNotFound {
		// create check
		err, table = CreateTable(p, pg)
		if err == nil {
			table.CreateCheck()
		}

	} else {
		log.Logger.Errorf("CreateRecord err: %s", err)
		return err, table
	}
	return nil, table
}
func CreateTable(p *restaurantmodels.Party, pg *restaurantmodels.PartyGroup) (err error, table models.Table) {
	var rTable restaurantmodels.Table
	var gOrder gatewaymodels.Order
	if err = db.RestDB.Model(&rTable).Find(&rTable, p.TableID.V()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
		}
		log.Logger.Errorf("CreateTable find table err: %s", err)
		return
	}
	if err = db.GatewayDB.Model(&gOrder).Find(&gOrder, pg.OrderID.V()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
		}
		log.Logger.Errorf("CreateTable find order err: %s", err)
		return
	}
	table.Table = rTable.Name.V()
	table.RestaurantTableID = rTable.ID.V()
	table.RestaurantPartyID = p.ID.V()
	table.RestaurantUpdatedAt = &p.UpdatedAt
	table.GuestCount = p.People.V()
	if p.SeatedAt != nil {
		table.CheckOpenTime = p.SeatedAt.Unix()
	}
	if p.UnseatedAt != nil {
		table.CheckCloseTime = p.UnseatedAt.Unix()
	}
	table.CheckRef = gOrder.Number.V()
	table.OrderID = gOrder.ID.V()
	table.StoreID = pg.StoreID.V()
	table.Status = models.TableStatusCreate
	table.HasSync = false
	err = db.DB.Model(&table).Create(&table).Error
	if err != nil {
		log.Logger.Errorf("PartyID %d Create table err: %s", p.ID.V(), err)
	}
	return
}

func CloseChecks() {
	var table models.Table
	var tables []models.Table
	var err error
	if err = db.DB.Model(&table).Where("`has_sync` = ? and `status` = ?", false, models.TableStatusClose).Find(&tables).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("NewChecks err: %s", err)
		return
	}
	fmt.Printf("\n %d CloseChecks", len(tables))
	for _, t := range tables {
		t.CloseCheck()
	}
	return
}

func CreateChecks() {
	var tables []models.Table
	var table models.Table
	var err error
	if err = db.DB.Model(&table).Where("has_sync = ? and status = ?", false, models.TableStatusCreate).Find(&tables).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("NewChecks err: %s", err)
		return
	}
	fmt.Printf("\n %d CreateChecks", len(tables))
	for _, t := range tables {
		t.CreateCheck()
	}
	return
}

func UpdateChecks() {
	var tables []models.Table
	var table models.Table
	var err error
	if err = db.DB.Model(&table).Where("has_sync = ? and status = ?", false, models.TableStatusUpdate).Find(&tables).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("NewChecks err: %s", err)
		return
	}
	fmt.Printf("\n %d UpdateChecks", len(tables))
	for _, t := range tables {
		t.UpdateCheck()
	}
	return
}

func MoveTables() {
	var tables []models.Table
	var table models.Table
	var err error
	if err = db.DB.Model(&table).Where("has_sync = ? and status = ?", false, models.TableStatusMove).Find(&tables).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("NewChecks err: %s", err)
		return
	}
	fmt.Printf("\n %d MoveTables", len(tables))
	for _, t := range tables {
		t.MoveTable()
	}
	return
}

func VoidChecks() {
	var tables []models.Table
	var table models.Table
	var err error
	if err = db.DB.Model(&table).Where("has_sync = ? and status = ?", false, models.TableStatusVoid).Find(&tables).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("NewChecks err: %s", err)
		return
	}
	fmt.Printf("\n %d UpdateChecks", len(tables))
	for _, t := range tables {
		t.VoidCheck()
	}
	return
}

type Resp struct {
	Result string `json:"result"`
}

func GetLastUpdatePartyTime() (time.Time, error) {
	var table models.Table
	ct := time.Now().Add(-24 * time.Hour)
	err := db.DB.Model(&table).Select("restaurant_updated_at").Order("restaurant_updated_at desc").Limit(1).Find(&table).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ct, nil
		}
		return ct, err
	}

	return *table.RestaurantUpdatedAt, nil
}

func ScaneParties() (err error, rParties []restaurantmodels.Party) {
	dt, err := GetLastUpdatePartyTime()
	fmt.Printf("GetLastUpdatePartyTime %+v", dt)
	du := time.Duration(-10 * 60 * 1000 * 1000 * 1000) //允许10分钟的通讯延迟
	dt = dt.Add(du)
	var rParty restaurantmodels.Party
	err = db.RestDB.Model(&rParty).Where("`updated_at` >= ?", dt).Order("updated_at asc").Find(&rParties).Error
	// err = db.RestDB.Model(&rParty).Where("`updated_at` > ?", "2019-01-29 00:04:00").Order("updated_at asc").Find(&rParties).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		log.Logger.Errorf("ScaneParties err: %s", err)
	}
	fmt.Printf("\n we find %d rParties", len(rParties))
	return
}

func FilterPartyBySettingStores(rParties []restaurantmodels.Party) {
	for _, p := range rParties {
		var rPartyGroup restaurantmodels.PartyGroup
		err := db.RestDB.Model(&rPartyGroup).Find(&rPartyGroup, p.PartyGroupID.V()).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			log.Logger.Errorf("FilterPartyBySettingStores err: %s", err)
			return
		}
		var isInSettingStore bool
		for _, store := range config.Conf.Setting.Stores {
			if rPartyGroup.StoreID.V() == store.StoreID {
				isInSettingStore = true
			}
		}
		if isInSettingStore {
			fmt.Printf("\n Analyzing party group %d", rPartyGroup.ID.V())
			AnalyzeRecord(&p, &rPartyGroup)
		}
	}
}

func ListingVoidPayments() {
	var tables []models.Table
	var table models.Table
	var err error
	ct := time.Now().Add(-8 * time.Hour)
	if err = db.DB.Model(&table).Where("has_sync = ? and created_at > ?", true, ct).Find(&tables).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("NewChecks err: %s", err)
		return
	}
	fmt.Printf("\n ListingVoidPayments get %d tables", len(tables))
	for _, t := range tables {
		t.GetAndSaveTransactions()
	}
}

func VoidPaymentOrCloseChecks() {
	var trans []models.Transaction
	var tran models.Transaction
	var err error
	if err = db.DB.Model(&tran).Where("has_sync = ?", false).Order("gateway_created_at asc").Find(&trans).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		log.Logger.Errorf("VoidOrCloseChecks err: %s", err)
		return
	}
	fmt.Printf("\n %d VoidOrCloseChecks", len(trans))
	for _, t := range trans {
		t.VoidPaymentOrCloseCheck()
	}
}
