package models

import (
	"bindolabs/optitable_middleware/config"
	"bindolabs/optitable_middleware/db"
	"bindolabs/optitable_middleware/gatewaymodels"
	"bindolabs/optitable_middleware/log"
	"bindolabs/optitable_middleware/optitable"
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
)

type Table struct {
	gorm.Model
	RestaurantPartyID   int        `gorm:"column:restaurant_party_id; type:int(11) ;" json:"restaurant_group_id"`
	RestaurantTableID   int        `gorm:"column:restaurant_table_id; type:int(11) ;" json:"restaurant_table_id"`
	OrderID             int        `gorm:"column:order_id; type:int(11) ;" json:"order_id"`
	RestaurantUpdatedAt *time.Time `gorm:"column:restaurant_updated_at; timestamp DEFAULT NULL ;" json:"restaurant_updated_at"`
	Status              int        `gorm:"column:status; type:int(11) ;" json:"status"`
	StoreID             int        `gorm:"column:store_id; type:int(11) ;" json:"store_id"`
	CheckRef            string     `gorm:"column:check_ref; type:varchar(255) ;" json:"check_ref"`
	CheckOpenTime       int64      `gorm:"column:check_open_time; type:int(64) ;" json:"check_open_time"`
	CheckCloseTime      int64      `gorm:"column:check_close_time; type:int(64) ;" json:"check_close_time"`
	Table               string     `gorm:"column:table; type:varchar(255) ;" json:"table"`
	OrgTable            string     `gorm:"column:org_table; type:varchar(255) ;" json:"org_table"`
	TotalAmount         float64    `gorm:"column:total_amount; type:decimal(12,2) ;" json:"total_amount"`
	GuestCount          int        `gorm:"column:guest_count; type:int(11) ;" json:"guest_count"`
	ChildCount          int        `gorm:"column:chind_count; type:int(11) ;" json:"chind_count"`
	HasSync             bool       `gorm:"column:has_sync; type:tinyint(1) ;default:'0'" json:"has_sync"`
}

const (
	TableStatusCreate = 1
	TableStatusClose  = 2
	TableStatusUpdate = 3
	TableStatusMove   = 4
	TableStatusVoid   = 5
	TableStatusCancle = 6
)

func (*Table) TableName() string {
	return "tables"
}

type Resp struct {
	Result string `json:"result"`
}

func (table *Table) UpdateCheck() (err error) {
	var resp Resp
	param := make(url.Values, 7)
	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	param.Set("check_ref", table.CheckRef)
	param.Set("check_open_time", fmt.Sprintf("%d", table.CheckOpenTime))
	param.Set("check_close_time", fmt.Sprintf("%d", table.CheckCloseTime))
	param.Set("work_date", table.RestaurantUpdatedAt.Format("2006-01-02"))
	param.Set("table", table.Table)
	if table.GuestCount == 0 {
		param.Set("guest_count", "1")
	} else {
		param.Set("guest_count", fmt.Sprintf("%d", table.GuestCount))
	}
	param.Set("child_count", fmt.Sprintf("%d", table.ChildCount))
	param.Set("api-key", store.OpApiKey)
	param.Set("function", "update_check")
	// param.Set("function", "close_check")

	err = optitable.Get(&param, &resp)
	log.Logger.Infof("update_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get optitable err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get optitable err, resp: %+v", resp)
		log.Logger.Error(err)
		return
	}
	table.HasSync = true
	table.Status = TableStatusUpdate
	err = db.DB.Model(&table).Save(&table).Error
	if err != nil {
		log.Logger.Errorf("UpdateCheck save table err: %s", err)
	}
	return
}

func (table *Table) CloseCheck() (err error) {
	var resp Resp
	param := make(url.Values, 7)
	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}

	param.Set("check_ref", table.CheckRef)
	param.Set("check_open_time", fmt.Sprintf("%d", table.CheckOpenTime))
	param.Set("check_close_time", fmt.Sprintf("%d", time.Now().Unix()))
	param.Set("table", table.Table)
	if table.GuestCount == 0 {
		param.Set("guest_count", "1")
	} else {
		param.Set("guest_count", fmt.Sprintf("%d", table.GuestCount))
	}
	param.Set("child_count", fmt.Sprintf("%d", table.ChildCount))
	param.Set("total_amount", fmt.Sprintf("%f", table.TotalAmount))
	param.Set("api-key", store.OpApiKey)
	param.Set("function", "close_check")

	err = optitable.Get(&param, &resp)
	log.Logger.Infof("close_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get optitable err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get optitable err, resp: %+v", resp)
		log.Logger.Error(err)
		return
	}
	table.HasSync = true
	table.Status = TableStatusClose
	err = db.DB.Model(&table).Save(&table).Error
	if err != nil {
		log.Logger.Errorf("CloseChecks save table err: %s", err)
	}
	return
}

func (table *Table) CreateCheck() (err error) {
	var resp Resp
	param := make(url.Values, 7)

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}

	param.Set("check_ref", table.CheckRef)
	param.Set("check_open_time", fmt.Sprintf("%d", table.CheckOpenTime))
	param.Set("table", table.Table)
	// param.Set("table", fmt.Sprintf("T%s", table.Table))
	if table.GuestCount == 0 {
		param.Set("guest_count", "1")
	} else {
		param.Set("guest_count", fmt.Sprintf("%d", table.GuestCount))
	}
	param.Set("child_count", fmt.Sprintf("%d", table.ChildCount))
	param.Set("api-key", store.OpApiKey)
	param.Set("function", "new_check")
	err = optitable.Get(&param, &resp)
	log.Logger.Infof("new_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get optitable err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get optitable err, resp: %+v", resp)
		log.Logger.Error(err)
		return
	}
	table.CheckCloseTime = 0
	table.HasSync = true
	table.Status = TableStatusCreate
	err = db.DB.Model(&table).Save(&table).Error
	if err != nil {
		log.Logger.Errorf("NewChecks save table err: %s", err)
	}
	return
}

func (table *Table) MoveTable() (err error) {

	var resp Resp
	param := make(url.Values, 7)

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}

	param.Set("check_ref", table.CheckRef)
	param.Set("check_open_time", fmt.Sprintf("%d", table.CheckOpenTime))
	param.Set("target_table", table.Table)
	param.Set("org_table", table.OrgTable)
	param.Set("work_date", table.RestaurantUpdatedAt.Format("2006-01-02"))
	if table.GuestCount == 0 {
		param.Set("guest_count", "1")
	} else {
		param.Set("guest_count", fmt.Sprintf("%d", table.GuestCount))
	}

	param.Set("child_count", fmt.Sprintf("%d", table.ChildCount))
	param.Set("api-key", store.OpApiKey)
	param.Set("function", "move_table")
	err = optitable.Get(&param, &resp)
	log.Logger.Infof("new_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get optitable err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get optitable err, resp: %+v", resp)
		log.Logger.Error(err)
		return
	}

	table.HasSync = true
	table.Status = TableStatusMove
	err = db.DB.Model(&table).Save(&table).Error
	if err != nil {
		log.Logger.Errorf("NewChecks save table err: %s", err)
	}
	return
}

func (table *Table) VoidCheck() (err error) {
	var resp Resp
	param := make(url.Values, 7)

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}

	param.Set("check_ref", table.CheckRef)
	param.Set("check_open_time", fmt.Sprintf("%d", table.CheckOpenTime))
	param.Set("target_table", table.Table)
	param.Set("work_date", table.RestaurantUpdatedAt.Format("2006-01-02"))
	if table.GuestCount == 0 {
		param.Set("guest_count", "1")
	} else {
		param.Set("guest_count", fmt.Sprintf("%d", table.GuestCount))
	}

	param.Set("child_count", fmt.Sprintf("%d", table.ChildCount))
	param.Set("api-key", store.OpApiKey)
	param.Set("function", "void_check")
	err = optitable.Get(&param, &resp)
	log.Logger.Infof("new_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get optitable err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get optitable err, resp: %+v", resp)
		log.Logger.Error(err)
		return
	}

	table.HasSync = true
	table.Status = TableStatusVoid
	err = db.DB.Model(&table).Save(&table).Error
	if err != nil {
		log.Logger.Errorf("VoidCheck save table err: %s", err)
	}
	return
}

func (table *Table) GetTransactions() (err error, trans []gatewaymodels.Transaction) {
	var tran gatewaymodels.Transaction
	if table.OrderID == 0 {
		err = fmt.Errorf("table.OrderID is 0")
		return
	}
	if err = db.GatewayDB.Model(&tran).Where("`source_type` = 'Order' and `source_id` = ?", table.OrderID).Find(&trans).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
		}
		log.Logger.Errorf("GetTransaction find Transaction err: %s", err)
		return
	}
	return
}

func (table *Table) GetAndSaveTransactions() (err error, trans []Transaction) {
	err, gwtrans := table.GetTransactions()
	fmt.Printf("\n trans(%d)", len(trans))
	if err != nil {
		log.Logger.Errorf("GetOrder err: %s", err)
		return
	}
	fmt.Printf("\n table %d GetTransactions %d \n", table.ID, len(trans))
	for _, gwt := range gwtrans {
		var tran Transaction
		if err = db.DB.Model(&tran).Where("gateway_transaction_id = ?", gwt.ID.V()).Find(&tran).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var newTran Transaction
				newTran.GatewayTransactionID = gwt.ID.V()
				newTran.GatewayCreatedAt = gwt.CreatedAt
				newTran.GatewayAction = gwt.Action.V()
				newTran.GatewayVoidedAt = gwt.VoidedAt
				newTran.Amount = gwt.Amount.V()
				newTran.TableID = table.ID
				err = db.DB.Model(&newTran).Save(&newTran).Error
				if err != nil {
					log.Logger.Errorf(" save Transaction err: %s", err)
				}
				return
			}
			log.Logger.Errorf("ListingVoidPayments err: %s", err)
			return
		}
		var tranVoidedAt int64
		var gwVoidedAt int64
		if tran.GatewayVoidedAt != nil {
			tranVoidedAt = tran.GatewayVoidedAt.Unix()
		}
		if gwt.VoidedAt != nil {
			gwVoidedAt = gwt.VoidedAt.Unix()
		}
		if tranVoidedAt != gwVoidedAt {
			tran.GatewayVoidedAt = gwt.VoidedAt
			tran.HasSync = false
			err = db.DB.Model(&tran).Save(&tran).Error
			if err != nil {
				log.Logger.Errorf(" save Transaction err: %s", err)
			}
		}
		if tran.HasSync == false {
			trans = append(trans, tran)
		}
	}
	return
}

func (table *Table) CanclePayment() (err error) {
	var resp Resp
	param := make(url.Values, 7)

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}

	param.Set("check_ref", table.CheckRef)
	param.Set("check_open_time", fmt.Sprintf("%d", table.CheckOpenTime))
	param.Set("target_table", table.Table)
	param.Set("work_date", table.RestaurantUpdatedAt.Format("2006-01-02"))
	if table.GuestCount == 0 {
		param.Set("guest_count", "1")
	} else {
		param.Set("guest_count", fmt.Sprintf("%d", table.GuestCount))
	}

	param.Set("child_count", fmt.Sprintf("%d", table.ChildCount))
	param.Set("api-key", store.OpApiKey)
	param.Set("function", "cancel_payment")
	err = optitable.Get(&param, &resp)
	log.Logger.Infof("cancle_payment response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("CanclePayment err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("CanclePayment err, resp: %+v", resp)
		log.Logger.Error(err)
		return
	}

	table.HasSync = true
	table.Status = TableStatusCancle
	err = db.DB.Model(&table).Save(&table).Error
	if err != nil {
		log.Logger.Errorf("VoidCheck save table err: %s", err)
	}
	return
}
