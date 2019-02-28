package models

import (
	"bindolabs/sevenrooms/config"
	"bindolabs/sevenrooms/db"
	"bindolabs/sevenrooms/gatewaymodels"
	"bindolabs/sevenrooms/log"

	"bindolabs/sevenrooms/sevenroom"
	"fmt"
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
	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	params := map[string]interface{}{
		"function":         "update_cover",
		"check_open_time":  time.Unix(table.CheckOpenTime, 0).Format(time.RFC3339),
		"table":            table.Table,
		"order_number":     table.CheckRef,
		"store_id":         store.StoreID,
		"store_name":       store.Name,
		"check_close_time": time.Now().Format(time.RFC3339),
		"total_amount":     table.TotalAmount,
	}
	if table.GuestCount == 0 {
		params["guest_count"] = 1
	} else {
		params["guest_count"] = table.GuestCount
	}
	err = sevenroom.PostWebhooks(store.VenueID, &params, &resp)

	log.Logger.Infof("update_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get sevenroom err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get sevenroom err, resp: %+v", resp)
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
	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	err, lineItems := table.GetLineItems()
	if err != nil {
		log.Logger.Errorf("LineItems err: %s", err)
		return
	}
	var lineItemData []map[string]interface{}
	for _, l := range lineItems {
		a := map[string]interface{}{
			"listing_barcode": l.ListingBarcode,
			"name":            l.Name,
			"qty":             l.Qty,
			"price":           l.Price,
			"net_total":       l.NetTotal,
		}
		lineItemData = append(lineItemData, a)
	}
	params := map[string]interface{}{
		"function":         "close_check",
		"check_open_time":  time.Unix(table.CheckOpenTime, 0).Format(time.RFC3339),
		"table":            table.Table,
		"order_number":     table.CheckRef,
		"store_id":         store.StoreID,
		"store_name":       store.Name,
		"check_close_time": time.Now().Format(time.RFC3339),
		"total_amount":     table.TotalAmount,
		"item":             lineItemData,
	}
	if table.GuestCount == 0 {
		params["guest_count"] = 1
	} else {
		params["guest_count"] = table.GuestCount
	}
	err = sevenroom.PostWebhooks(store.VenueID, &params, &resp)
	log.Logger.Infof("close_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get sevenroom err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get sevenroom err, resp: %+v", resp)
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

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	params := map[string]interface{}{
		"function":        "new_check",
		"check_open_time": time.Unix(table.CheckOpenTime, 0).Format(time.RFC3339),
		"table":           table.Table,
		"order_number":    table.CheckRef,
		"store_id":        store.StoreID,
		"store_name":      store.Name,
	}
	if table.GuestCount == 0 {
		params["guest_count"] = 1
	} else {
		params["guest_count"] = table.GuestCount
	}
	err = sevenroom.PostWebhooks(store.VenueID, &params, &resp)
	log.Logger.Infof("new_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("CreateCheck err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get sevenroom err, resp: %+v", resp)
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

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	params := map[string]interface{}{
		"function":        "move_table",
		"check_open_time": time.Unix(table.CheckOpenTime, 0).Format(time.RFC3339),
		"table":           table.Table,
		"order_number":    table.CheckRef,
		"store_id":        store.StoreID,
		"store_name":      store.Name,
		"from_table":      table.OrgTable,
		"to_table":        table.Table,
	}
	if table.GuestCount == 0 {
		params["guest_count"] = 1
	} else {
		params["guest_count"] = table.GuestCount
	}
	err = sevenroom.PostWebhooks(store.VenueID, &params, &resp)

	log.Logger.Infof("new_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get sevenroom err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get sevenroom err, resp: %+v", resp)
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

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	params := map[string]interface{}{
		"function":        "void_check",
		"check_open_time": time.Unix(table.CheckOpenTime, 0).Format(time.RFC3339),
		"table":           table.Table,
		"order_number":    table.CheckRef,
		"store_id":        store.StoreID,
		"store_name":      store.Name,
	}
	if table.GuestCount == 0 {
		params["guest_count"] = 1
	} else {
		params["guest_count"] = table.GuestCount
	}
	err = sevenroom.PostWebhooks(store.VenueID, &params, &resp)

	log.Logger.Infof("new_check response: %+v", resp)
	if err != nil {
		log.Logger.Errorf("Get sevenroom err %s", err)
		return
	}
	if resp.Result != "success" {
		err = fmt.Errorf("Get sevenroom err, resp: %+v", resp)
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

	err, store := config.GetStore(table.StoreID)
	if err != nil {
		log.Logger.Errorf("GetStore err: %s", err)
		return
	}
	params := map[string]interface{}{
		"function":        "cancel_payment",
		"check_open_time": time.Unix(table.CheckOpenTime, 0).Format(time.RFC3339),
		"table":           table.Table,
		"order_number":    table.CheckRef,
		"store_id":        store.StoreID,
		"store_name":      store.Name,
	}
	if table.GuestCount == 0 {
		params["guest_count"] = 1
	} else {
		params["guest_count"] = table.GuestCount
	}
	err = sevenroom.PostWebhooks(store.VenueID, &params, &resp)
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

func (table *Table) GetLineItems() (err error, lineItems []LineItem) {
	var gLineItem gatewaymodels.LineItem
	var gLineItems []gatewaymodels.LineItem
	var lineItem LineItem
	if err = db.GatewayDB.Model(&gLineItem).Where("order_id = ?", table.OrderID).Find(&gLineItems).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
		}
		log.Logger.Errorf("find GatewayDB LineItems err: %s", err)
		return
	}
	if err = db.DB.Model(&lineItem).Where("table_id = ? ", table.ID).Delete("").Error; err != nil {
		log.Logger.Errorf("delete LineItems err: %s", err)
		return
	}

	for _, l := range gLineItems {
		var lineItem LineItem
		lineItem.TableID = table.ID
		lineItem.GatewayLineItemID = l.ID.V()
		lineItem.Name = l.Label.V()
		lineItem.Qty = l.Quantity.V()
		lineItem.Price = l.Price.V()
		lineItem.NetTotal = l.Total.V()
		lineItem.GetListingBarcode()
		if err = db.DB.Model(&lineItem).Create(&lineItem).Error; err != nil {
			log.Logger.Errorf("create LineItems err: %s", err)
			return
		}
		lineItems = append(lineItems, lineItem)

	}
	return
}
