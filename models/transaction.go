package models

import (
	"bindolabs/sevenrooms/db"
	"bindolabs/sevenrooms/log"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	Amount               float64    `gorm:"column:amount; type:decimal(16,2) ;default:'0.00'" json:"amount"`
	GatewayTransactionID int        `gorm:"column:gateway_transaction_id; type:int(11) ;" json:"gateway_transaction_id"`
	GatewayCreatedAt     time.Time  `gorm:"column:gateway_created_at; type:datetime ;" json:"gateway_created_at"`
	GatewayAction        int        `gorm:"column:gateway_action; type:int(11) ;" json:"gateway_action"`
	GatewayVoidedAt      *time.Time `gorm:"column:gateway_voided_at; type:datetime ;" json:"gateway_voided_at"`
	TableID              uint       `gorm:"column:table_id; type:int(11) ;" json:"table_id"`
	HasSync              bool       `gorm:"column:has_sync; type:tinyint(1) ;default:'0'" json:"has_sync"`
}

func (*Transaction) TableName() string {
	return "transactions"
}

const (
	TransactionActionRreturn = 1 // refund transaction, source_type === 'Refund'
	TransactionActionSale    = 2 // any payment, include both partial paid & full paid
	TransactionActionToPUP   = 3
	TransactionActionAuth    = 4
	TransactionActionAdjTips = 5  // add tips
	TransactionActionSign    = 8  // for credit card signiture
	TransactionActionVoid    = 16 // for credit card, source_type === 'Refund'
)

func (t *Transaction) VoidPaymentOrCloseCheck() {
	fmt.Printf("\n tran %+v", t)
	fmt.Printf("\n t.GatewayVoidedAt %+v", t.GatewayVoidedAt)
	fmt.Printf("\n t.GatewayVoidedAt != nil %+v", t.GatewayVoidedAt != nil)
	fmt.Printf("\n t.GatewayAction == models.TransactionActionSale %+v", t.GatewayAction == TransactionActionSale)
	if t.GatewayAction == TransactionActionSale {
		var table Table
		var err error
		if err = db.DB.Model(&table).Where("has_sync = ?", true).Find(&table, t.TableID).Error; err != nil {
			log.Logger.Errorf("VoidOrCloseChecks find table err: %s", err)
			return
		}
		if t.GatewayVoidedAt == nil && table.Status != TableStatusClose {
			table.TotalAmount = t.Amount
			err = db.DB.Model(&table).Save(&table).Error
			if err != nil {
				log.Logger.Errorf("VoidPaymentOrCloseChecks save table err: %s", err)
				return
			}
			err = table.CloseCheck()
			if err == nil {
				t.HasSync = true
				err = db.DB.Model(&t).Save(&t).Error
				if err != nil {
					log.Logger.Errorf("VoidPaymentOrCloseChecks save Transaction err: %s", err)
					return
				}
			}
		} else if t.GatewayVoidedAt != nil {
			fmt.Println("\n +++++++CanclePayment")
			err = table.CanclePayment()
			if err == nil {
				t.HasSync = true
				err = db.DB.Model(&t).Save(&t).Error
				if err != nil {
					log.Logger.Errorf("CanclePayment save Transaction err: %s", err)
					return
				}
			}
			// var trans2 []models.Transaction
			// if err = db.DB.Model(&tran).Where("has_sync = ?", false).Order("gateway_created_at asc").Find(&trans2).Error; err != nil {
			// 	if err == gorm.ErrRecordNotFound {
			// 		return
			// 	}
			// 	log.Logger.Errorf("VoidOrCloseChecks err: %s", err)
			// 	return
			// }
			// fmt.Printf("\n %d ========trans2", len(trans2))
			// table.CreateCheck()
		}
	}
}
