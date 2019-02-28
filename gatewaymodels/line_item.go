package gatewaymodels

import (
	"time"

	"github.com/imiskolee/optional"
)

type LineItem struct {
	ID                      optional.Int     `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	AddOnTaxAmount          optional.Float64 `gorm:"column:add_on_tax_amount; type:decimal(16,2) ;" json:"add_on_tax_amount"`
	BaseUnit                optional.String  `gorm:"column:base_unit; type:varchar(255) ;" json:"base_unit"`
	BaseUnitID              optional.Int     `gorm:"column:base_unit_id; type:int(11) ;" json:"base_unit_id"`
	CashierID               optional.Int     `gorm:"column:cashier_id; type:int(11) ;" json:"cashier_id"`
	CourierID               optional.Int     `gorm:"column:courier_id; type:int(11) ;" json:"courier_id"`
	CreatedAt               *time.Time       `gorm:"column:created_at; type:datetime ;" json:"created_at"`
	DeletedAt               *time.Time       `gorm:"column:deleted_at; type:datetime ;" json:"deleted_at"`
	Description             optional.String  `gorm:"column:description; type:text ;" json:"description"`
	DiscountTotal           optional.Float64 `gorm:"column:discount_total; type:decimal(16,2) ;" json:"discount_total"`
	EffectiveCreatedAt      *time.Time       `gorm:"column:effective_created_at; type:datetime ;" json:"effective_created_at"`
	EffectiveDeletedAt      *time.Time       `gorm:"column:effective_deleted_at; type:datetime ;" json:"effective_deleted_at"`
	FavoriteID              optional.Int     `gorm:"column:favorite_id; type:int(11) ;" json:"favorite_id"`
	FavoriteSectionID       optional.Int     `gorm:"column:favorite_section_id; type:int(11) ;" json:"favorite_section_id"`
	FavoriteTabID           optional.Int     `gorm:"column:favorite_tab_id; type:int(11) ;" json:"favorite_tab_id"`
	FulfilmentNote          optional.String  `gorm:"column:fulfilment_note; type:mediumtext ;" json:"fulfilment_note"`
	GiftCardID              optional.Int     `gorm:"column:gift_card_id; type:int(11) ;" json:"gift_card_id"`
	GroupID                 optional.Int     `gorm:"column:group_id; type:int(11) ;" json:"group_id"`
	IdentifierForVendor     optional.String  `gorm:"column:identifier_for_vendor; type:varchar(255) ;" json:"identifier_for_vendor"`
	IncludeTaxPrice         optional.Float64 `gorm:"column:include_tax_price; type:decimal(16,2) ;" json:"include_tax_price"`
	InclusiveTaxAmount      optional.Float64 `gorm:"column:inclusive_tax_amount; type:decimal(16,2) ;" json:"inclusive_tax_amount"`
	KitchenNote             optional.String  `gorm:"column:kitchen_note; type:text ;" json:"kitchen_note"`
	Label                   optional.String  `gorm:"column:label; type:varchar(255) ;" json:"label"`
	LineItemCode            optional.String  `gorm:"column:line_item_code; type:varchar(255) ;" json:"line_item_code"`
	LineItemCodeActive      optional.Bool    `gorm:"column:line_item_code_active; type:tinyint(1) ;" json:"line_item_code_active"`
	ModifierSetOptionID     optional.Int     `gorm:"column:modifier_set_option_id; type:int(11) ;" json:"modifier_set_option_id"`
	Note                    optional.String  `gorm:"column:note; type:mediumtext ;" json:"note"`
	OrderID                 optional.Int     `gorm:"column:order_id; type:int(11) ;" json:"order_id"`
	OriginalPrice           optional.Float64 `gorm:"column:original_price; type:decimal(9,2) ;" json:"original_price"`
	Price                   optional.Float64 `gorm:"column:price; type:decimal(16,2) ;" json:"price"`
	PriceOptionID           optional.Int     `gorm:"column:price_option_id; type:int(11) ;" json:"price_option_id"`
	PriceOptionName         optional.String  `gorm:"column:price_option_name; type:varchar(255) ;" json:"price_option_name"`
	PricePerUnit            optional.Float64 `gorm:"column:price_per_unit; type:decimal(16,2) ;" json:"price_per_unit"`
	PurchasableID           optional.Int     `gorm:"column:purchasable_id; type:int(11) ;" json:"purchasable_id"`
	PurchasableType         optional.String  `gorm:"column:purchasable_type; type:varchar(255) ;" json:"purchasable_type"`
	Quantity                optional.Float64 `gorm:"column:quantity; type:decimal(12,4) ;" json:"quantity"`
	QuantityAllowDecimal    optional.Bool    `gorm:"column:quantity_allow_decimal; type:tinyint(1) ;default:'0'" json:"quantity_allow_decimal"`
	RoundingAmount          optional.Float64 `gorm:"column:rounding_amount; type:decimal(16,2) ;default:'0.00'" json:"rounding_amount"`
	SentToKitchenCount      optional.Int     `gorm:"column:sent_to_kitchen_count; type:int(11) ;default:'0'" json:"sent_to_kitchen_count"`
	ShippingMethod          optional.Int     `gorm:"column:shipping_method; type:int(11) ;" json:"shipping_method"`
	SplitError              optional.Float64 `gorm:"column:split_error; type:decimal(16,2) ;default:'0.00'" json:"split_error"`
	StoreID                 optional.Int     `gorm:"column:store_id; type:int(11) ;" json:"store_id"`
	Subtotal                optional.Float64 `gorm:"column:subtotal; type:decimal(16,2) ;" json:"subtotal"`
	TableSeat               optional.String  `gorm:"column:table_seat; type:varchar(255) ;" json:"table_seat"`
	TaxAmount               optional.Float64 `gorm:"column:tax_amount; type:decimal(16,7) ;" json:"tax_amount"`
	TaxBaseAdjustmentAmount optional.Float64 `gorm:"column:tax_base_adjustment_amount; type:decimal(16,2) ;" json:"tax_base_adjustment_amount"`
	TaxInPriceRoundingError optional.Float64 `gorm:"column:tax_in_price_rounding_error; type:decimal(16,2) ;default:'0.00'" json:"tax_in_price_rounding_error"`
	TaxOptionID             optional.Int     `gorm:"column:tax_option_id; type:int(11) ;" json:"tax_option_id"`
	TaxRate                 optional.Float64 `gorm:"column:tax_rate; type:decimal(10,7) ;" json:"tax_rate"`
	TaxableAmount           optional.Float64 `gorm:"column:taxable_amount; type:decimal(16,2) ;default:'0.00'" json:"taxable_amount"`
	Total                   optional.Float64 `gorm:"column:total; type:decimal(16,2) ;" json:"total"`
	Unit                    optional.String  `gorm:"column:unit; type:varchar(255) ;" json:"unit"`
	UnitID                  optional.Int     `gorm:"column:unit_id; type:int(11) ;" json:"unit_id"`
	UnitQuantity            optional.Float64 `gorm:"column:unit_quantity; type:decimal(16,4) ;" json:"unit_quantity"`
	UpdatedAt               *time.Time       `gorm:"column:updated_at; type:datetime ;" json:"updated_at"`
	UUID                    []byte           `gorm:"column:uuid; type:binary(16) ;" json:"uuid"`
	VoidApprovedBy          optional.Int     `gorm:"column:void_approved_by; type:int(11) ;" json:"void_approved_by"`
	VoidNote                optional.String  `gorm:"column:void_note; type:text ;" json:"void_note"`
	VoidReason              optional.String  `gorm:"column:void_reason; type:varchar(255) ;" json:"void_reason"`
	Voided                  optional.Bool    `gorm:"column:voided; type:tinyint(1) ;default:'0'" json:"voided"`
	VoidedBy                optional.Int     `gorm:"column:voided_by; type:int(11) ;" json:"voided_by"`
}

// TableName sets the insert table name for this struct type
func (l *LineItem) TableName() string {
	return "line_items"
}
