package gatewaymodels

import (
	"time"

	"github.com/imiskolee/optional"
)

type Transaction struct {
	AcquiringPartnerID             optional.Int     `gorm:"column:acquiring_partner_id; type:int(11) ;" json:"acquiring_partner_id"`
	Action                         optional.Int     `gorm:"column:action; type:int(11) ;" json:"action"`
	Amount                         optional.Float64 `gorm:"column:amount; type:decimal(16,2) ;default:'0.00'" json:"amount"`
	BatchRoutineID                 optional.Int     `gorm:"column:batch_routine_id; type:int(11) ;" json:"batch_routine_id"`
	BindoAcquiringBatchID          optional.Int     `gorm:"column:bindo_acquiring_batch_id; type:int(11) ;" json:"bindo_acquiring_batch_id"`
	BindoMerchantFeeBatchID        optional.Int     `gorm:"column:bindo_merchant_fee_batch_id; type:int(11) ;" json:"bindo_merchant_fee_batch_id"`
	BindoMerchantFloatDay          optional.Int     `gorm:"column:bindo_merchant_float_day; type:int(11) ;" json:"bindo_merchant_float_day"`
	BindoMerchantSaleBatchID       optional.Int     `gorm:"column:bindo_merchant_sale_batch_id; type:int(11) ;" json:"bindo_merchant_sale_batch_id"`
	BindoSpread                    optional.Float64 `gorm:"column:bindo_spread; type:decimal(16,2) ;" json:"bindo_spread"`
	BindoTotalFee                  optional.Float64 `gorm:"column:bindo_total_fee; type:decimal(16,2) ;" json:"bindo_total_fee"`
	Canceled                       optional.Bool    `gorm:"column:canceled; type:tinyint(1) ;default:'0'" json:"canceled"`
	CaptureAmount                  optional.Float64 `gorm:"column:capture_amount; type:decimal(16,2) ;" json:"capture_amount"`
	Captured                       optional.Bool    `gorm:"column:captured; type:tinyint(1) ;default:'0'" json:"captured"`
	CardReaderModel                optional.Int     `gorm:"column:card_reader_model; type:int(11) ;" json:"card_reader_model"`
	CashierID                      optional.Int     `gorm:"column:cashier_id; type:int(11) ;" json:"cashier_id"`
	CashierName                    optional.String  `gorm:"column:cashier_name; type:varchar(255) ;" json:"cashier_name"`
	ChangeAmount                   optional.Float64 `gorm:"column:change_amount; type:decimal(8,2) ;" json:"change_amount"`
	CreatedAt                      time.Time        `gorm:"column:created_at; type:datetime ;" json:"created_at"`
	CreditAuthDetail1RecordID      optional.Int     `gorm:"column:credit_auth_detail1_record_id; type:int(11) ;" json:"credit_auth_detail1_record_id"`
	CreditAuthDetail2RecordID      optional.Int     `gorm:"column:credit_auth_detail2_record_id; type:int(11) ;" json:"credit_auth_detail2_record_id"`
	CreditCardID                   optional.Int     `gorm:"column:credit_card_id; type:int(11) ;" json:"credit_card_id"`
	CreditReconDetail1RecordID     optional.Int     `gorm:"column:credit_recon_detail1_record_id; type:int(11) ;" json:"credit_recon_detail1_record_id"`
	CreditReconDetail2RecordID     optional.Int     `gorm:"column:credit_recon_detail2_record_id; type:int(11) ;" json:"credit_recon_detail2_record_id"`
	Currency                       optional.String  `gorm:"column:currency; type:varchar(255) ;default:'USD'" json:"currency"`
	EffectiveCreatedAt             *time.Time       `gorm:"column:effective_created_at; type:datetime ;" json:"effective_created_at"`
	FeeCcProcessing                optional.Float64 `gorm:"column:fee_cc_processing; type:decimal(8,2) ;default:'0.00'" json:"fee_cc_processing"`
	FeeCcRefund                    optional.Float64 `gorm:"column:fee_cc_refund; type:decimal(8,2) ;default:'0.00'" json:"fee_cc_refund"`
	FeeCcTransaction               optional.Float64 `gorm:"column:fee_cc_transaction; type:decimal(8,2) ;default:'0.00'" json:"fee_cc_transaction"`
	FeeCommission                  optional.Float64 `gorm:"column:fee_commission; type:decimal(8,2) ;default:'0.00'" json:"fee_commission"`
	From                           optional.Int     `gorm:"column:from; type:int(11) ;" json:"from"`
	ID                             optional.Int     `gorm:"column:id; type:int(11) AUTO_INCREMENT;" json:"id"`
	IdentifierForVendor            optional.String  `gorm:"column:identifier_for_vendor; type:varchar(255) ;" json:"identifier_for_vendor"`
	InterchangeFee                 optional.Float64 `gorm:"column:interchange_fee; type:decimal(16,2) ;" json:"interchange_fee"`
	InterchangePlus                optional.Bool    `gorm:"column:interchange_plus; type:tinyint(1) ;default:'0'" json:"interchange_plus"`
	MerchantAccountID              optional.Int     `gorm:"column:merchant_account_id; type:int(11) ;" json:"merchant_account_id"`
	Note                           optional.String  `gorm:"column:note; type:longtext ;" json:"note"`
	PaymentMethod                  optional.Int     `gorm:"column:payment_method; type:int(11) ;" json:"payment_method"`
	ReferenceNumber                optional.String  `gorm:"column:reference_number; type:varchar(255) ;" json:"reference_number"`
	ScheduleBindoMerchantBatchTime *time.Time       `gorm:"column:schedule_bindo_merchant_batch_time; type:datetime ;" json:"schedule_bindo_merchant_batch_time"`
	Settled                        optional.Bool    `gorm:"column:settled; type:tinyint(1) ;default:'0'" json:"settled"`
	SignatureRequired              optional.Bool    `gorm:"column:signature_required; type:tinyint(1) ;default:'0'" json:"signature_required"`
	SourceID                       optional.Int     `gorm:"column:source_id; type:int(11) ;" json:"source_id"`
	SourceType                     optional.String  `gorm:"column:source_type; type:varchar(255) ;" json:"source_type"`
	SpreadCcProcessing             optional.Float64 `gorm:"column:spread_cc_processing; type:decimal(8,2) ;default:'0.00'" json:"spread_cc_processing"`
	SpreadCcTransaction            optional.Float64 `gorm:"column:spread_cc_transaction; type:decimal(8,2) ;default:'0.00'" json:"spread_cc_transaction"`
	SpreadCommission               optional.Float64 `gorm:"column:spread_commission; type:decimal(8,2) ;default:'0.00'" json:"spread_commission"`
	StoreID                        optional.Int     `gorm:"column:store_id; type:int(11) ;default:'0'" json:"store_id"`
	Success                        optional.Bool    `gorm:"column:success; type:tinyint(1) ;default:'0'" json:"success"`
	IsPending                      optional.Bool    `gorm:"column:is_pending; type:tinyint(1) ;default:'0'" json:"is_pending"`
	TimeSegment                    optional.String  `gorm:"column:time_segment; type:varchar(255) ;" json:"time_segment"`
	TipsAmount                     optional.Float64 `gorm:"column:tips_amount; type:decimal(16,2) ;" json:"tips_amount"`
	TrailID                        optional.Int     `gorm:"column:trail_id; type:int(11) ;" json:"trail_id"`
	UUID                           []byte           `gorm:"column:uuid; type:binary(16) ;" json:"uuid"`
	ValidationErrors               optional.String  `gorm:"column:validation_errors; type:longtext ;" json:"validation_errors"`
	VoidedAt                       *time.Time       `gorm:"column:voided_at; type:datetime ;" json:"voided_at"`
	TransactionType                optional.String  `gorm:"column:transaction_type; type:CHAR(128) ;" json:"transaction_type"`
	IsReversal                     optional.Bool    `gorm:"column:is_reversal; type:tinyint(1) ;default:'0'" json:"is_reversal"`
	IsOffline                      optional.Bool    `gorm:"column:is_offline; type:tinyint(1) ;default:'0'" json:"is_offline"`
}

func (*Transaction) TableName() string {
	return "transactions"
}
