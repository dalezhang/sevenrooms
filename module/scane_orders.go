package module

import (
	"bindolabs/sevenrooms/db"
	"bindolabs/sevenrooms/gatewaymodels"
	"bindolabs/sevenrooms/log"
	"bindolabs/sevenrooms/models"
)

func ScaneOrders() {
	var table models.Table
	var tables []models.Table
	var err error
	err = db.DB.Model(&table).Where("status not in (2,5,6)").Find(&tables).Error
	if err != nil {
		log.Logger.Errorf("ScaneOrders find table err: %s", err)
	}
	for _, table := range tables {
		var order gatewaymodels.Order
		err = db.GatewayDB.Model(&order).Find(&order, table.OrderID).Error
		if err != nil {
			log.Logger.Errorf("ScaneOrders find order err: %s", err)
		}
		if order.Subtotal.V() != table.Subtotal {
			table.UpdateItems(order)
		}
	}
}
