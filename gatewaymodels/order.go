package gatewaymodels

import (
	"time"

	"github.com/imiskolee/optional"
)

type Order struct {
	AllowNegativeTotal                     optional.Bool    `gorm:"column:allow_negative_total; type:tinyint(1) ;default:'0'" json:"allow_negative_total"`
	AuthCode                               optional.String  `gorm:"column:auth_code; type:varchar(255) ;" json:"auth_code"`
	BarcodeBindoCodeContentType            optional.String  `gorm:"column:barcode_bindo_code_content_type; type:varchar(255) ;" json:"barcode_bindo_code_content_type"`
	BarcodeBindoCodeFileName               optional.String  `gorm:"column:barcode_bindo_code_file_name; type:varchar(255) ;" json:"barcode_bindo_code_file_name"`
	BarcodeBindoCodeFileSize               optional.Int     `gorm:"column:barcode_bindo_code_file_size; type:int(11) ;" json:"barcode_bindo_code_file_size"`
	BarcodeBindoCodeUpdatedAt              *time.Time       `gorm:"column:barcode_bindo_code_updated_at; type:datetime ;" json:"barcode_bindo_code_updated_at"`
	BarcodeNumberContentType               optional.String  `gorm:"column:barcode_number_content_type; type:varchar(255) ;" json:"barcode_number_content_type"`
	BarcodeNumberFileName                  optional.String  `gorm:"column:barcode_number_file_name; type:varchar(255) ;" json:"barcode_number_file_name"`
	BarcodeNumberFileSize                  optional.Int     `gorm:"column:barcode_number_file_size; type:int(11) ;" json:"barcode_number_file_size"`
	BarcodeNumberUpdatedAt                 *time.Time       `gorm:"column:barcode_number_updated_at; type:datetime ;" json:"barcode_number_updated_at"`
	BillingAddressID                       optional.Int     `gorm:"column:billing_address_id; type:int(11) ;" json:"billing_address_id"`
	BillingAddress                         optional.String  `gorm:"-" json:"billing_address"`
	CalculatorVersion                      optional.String  `gorm:"column:calculator_version; type:varchar(255) ;" json:"calculator_version"`
	CanceledAt                             *time.Time       `gorm:"column:canceled_at; type:datetime ;" json:"canceled_at"`
	CartID                                 optional.Int     `gorm:"column:cart_id; type:int(11) ;" json:"cart_id"`
	Checksum                               optional.String  `gorm:"column:checksum; type:varchar(255) ;" json:"checksum"`
	CodeConfirmedAt                        *time.Time       `gorm:"column:code_confirmed_at; type:datetime ;" json:"code_confirmed_at"`
	CompletedAt                            *time.Time       `gorm:"column:completed_at; type:datetime ;" json:"completed_at"`
	CorrespondenceState                    optional.String  `gorm:"column:correspondence_state; type:varchar(255) ;" json:"correspondence_state"`
	Counter                                optional.Int     `gorm:"column:counter; type:int(11) ;default:'0'" json:"counter"`
	CourierID                              optional.Int     `gorm:"column:courier_id; type:int(11) ;" json:"courier_id"`
	CreatedAt                              *time.Time       `gorm:"column:created_at; type:datetime ;" json:"created_at"`
	Currency                               optional.String  `gorm:"column:currency; type:varchar(255) ;default:'USD'" json:"currency"`
	CustomerID                             optional.Int     `gorm:"column:customer_id; type:int(11) ;" json:"customer_id"`
	DeliveredAt                            *time.Time       `gorm:"column:delivered_at; type:datetime ;" json:"delivered_at"`
	DeliveryDate                           *time.Time       `gorm:"column:delivery_date; type:datetime ;" json:"delivery_date"`
	DiscountTotal                          optional.Float64 `gorm:"column:discount_total; type:decimal(16,2) ;" json:"discount_total"`
	DueDate                                *time.Time       `gorm:"column:due_date; type:date ;" json:"due_date"`
	EffectiveCreatedAt                     *time.Time       `gorm:"column:effective_created_at; type:datetime ;" json:"effective_created_at"`
	Email                                  optional.String  `gorm:"column:email; type:varchar(255) ;" json:"email"`
	From                                   optional.Int     `gorm:"column:from; type:int(11) ;" json:"from"`
	ID                                     optional.Int     `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	IdentifierForVendor                    optional.String  `gorm:"column:identifier_for_vendor; type:varchar(255) ;" json:"identifier_for_vendor"`
	InitialDelivery                        optional.Float64 `gorm:"column:initial_delivery; type:decimal(16,2) ;" json:"initial_delivery"`
	InitialIncludedInPriceTax              optional.Float64 `gorm:"column:initial_included_in_price_tax; type:decimal(16,2) ;" json:"initial_included_in_price_tax"`
	InitialIncludedInPriceTaxError         optional.Float64 `gorm:"column:initial_included_in_price_tax_error; type:decimal(16,2) ;" json:"initial_included_in_price_tax_error"`
	InitialIncludedInPriceTaxForServiceFee optional.Float64 `gorm:"column:initial_included_in_price_tax_for_service_fee; type:decimal(16,2) ;" json:"initial_included_in_price_tax_for_service_fee"`
	InitialNormalTax                       optional.Float64 `gorm:"column:initial_normal_tax; type:decimal(16,2) ;" json:"initial_normal_tax"`
	InitialProductTotal                    optional.Float64 `gorm:"column:initial_product_total; type:decimal(16,2) ;" json:"initial_product_total"`
	InitialProvidedInRequest               optional.Bool    `gorm:"column:initial_provided_in_request; type:tinyint(1) ;default:'0'" json:"initial_provided_in_request"`
	InitialRedeemDeposits                  optional.Float64 `gorm:"column:initial_redeem_deposits; type:decimal(16,2) ;" json:"initial_redeem_deposits"`
	InitialRounding                        optional.Float64 `gorm:"column:initial_rounding; type:decimal(16,2) ;" json:"initial_rounding"`
	InitialServiceFee                      optional.Float64 `gorm:"column:initial_service_fee; type:decimal(16,2) ;" json:"initial_service_fee"`
	InitialTax                             optional.Float64 `gorm:"column:initial_tax; type:decimal(16,2) ;" json:"initial_tax"`
	InitialTips                            optional.Float64 `gorm:"column:initial_tips; type:decimal(16,2) ;default:'0.00'" json:"initial_tips"`
	InitialTotal                           optional.Float64 `gorm:"column:initial_total; type:decimal(16,2) ;" json:"initial_total"`
	IsReservation                          optional.Bool    `gorm:"column:is_reservation; type:tinyint(1) ;default:'0'" json:"is_reservation"`
	IsRegisterRefund                       optional.Bool    `gorm:"column:is_register_refund; type:tinyint(1) ;"  json:"is_register_refund"`
	Note                                   optional.String  `gorm:"column:note; type:mediumtext ;" json:"note"`
	Number                                 optional.String  `gorm:"column:number; type:varchar(255) ;" json:"number"`
	Offline                                optional.Bool    `gorm:"column:offline; type:tinyint(1) ;default:'0'" json:"offline"`
	PartialFulfilmentEnabled               optional.Bool    `gorm:"column:partial_fulfilment_enabled; type:tinyint(1) ;" json:"partial_fulfilment_enabled"`
	PickupLocationID                       optional.Int     `gorm:"column:pickup_location_id; type:int(11) ;" json:"pickup_location_id"`
	ReceiptCreatedAt                       *time.Time       `gorm:"column:receipt_created_at; type:datetime ;" json:"receipt_created_at"`
	ReceiptUpdatedAt                       *time.Time       `gorm:"column:receipt_updated_at; type:datetime ;" json:"receipt_updated_at"`
	ReferenceNumber                        optional.String  `gorm:"column:reference_number; type:varchar(255) ;" json:"reference_number"`
	ReturnType                             optional.Int     `gorm:"column:return_type; type:int(11) ;default:'0'" json:"return_type"`
	SaleType                               optional.String  `gorm:"column:sale_type; type:varchar(255) ;" json:"sale_type"`
	ShipToID                               optional.Int     `gorm:"column:ship_to_id; type:int(11) ;" json:"ship_to_id"`
	ShippingAddressID                      optional.Int     `gorm:"column:shipping_address_id; type:int(11) ;" json:"shipping_address_id"`
	ShippingAddress                        optional.String  `gorm:"-" json:"shipping_address"`
	ShippingMethod                         optional.Int     `gorm:"column:shipping_method; type:int(11) ;" json:"shipping_method"`
	State                                  optional.String  `gorm:"column:state; type:varchar(255) ;" json:"state"`
	StoreID                                optional.Int     `gorm:"column:store_id; type:int(11) ;" json:"store_id"`
	SubscriptionID                         optional.Int     `gorm:"column:subscription_id; type:int(11) ;" json:"subscription_id"`
	Subtotal                               optional.Float64 `gorm:"column:subtotal; type:decimal(16,2) ;" json:"subtotal"`
	TaxBaseAdjustment                      optional.Float64 `gorm:"column:tax_base_adjustment; type:decimal(16,2) ;" json:"tax_base_adjustment"`
	TimeSegment                            optional.String  `gorm:"column:time_segment; type:varchar(255) ;" json:"time_segment"`
	Transactionless                        optional.Bool    `gorm:"column:transactionless; type:tinyint(1) ;default:'0'" json:"transactionless"`
	UpdatedAt                              *time.Time       `gorm:"column:updated_at; type:datetime ;" json:"updated_at"`
	UserID                                 optional.Int     `gorm:"column:user_id; type:int(11) ;" json:"user_id"`
	UUID                                   []byte           `gorm:"column:uuid; type:binary(16) ;" json:"uuid"`
	Version                                optional.Int     `gorm:"column:version; type:int(11) ;" json:"version"`
	IsSuperOrder                           optional.Bool    `gorm:"column:is_super_order; type:tinyint(1) ;default:'0'"`
	PickupStatus                           optional.Int     `gorm:"column:pickup_status;type:int(11)"`
}

func (*Order) TableName() string {
	return "orders"
}
