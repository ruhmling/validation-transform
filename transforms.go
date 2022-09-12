package validationAndTransform

var SubsidizeTransaction Transform = func(payment Payment, data interface{}) interface{} {
	trx := data.(Transaction)
	trx.Amount = payment.Amount
	trx.StatusDetail = "amount_subsidized"

	return trx
}

var SetTransactionStatusCancelled Transform = func(payment Payment, data interface{}) interface{} {
	trx := data.(Transaction)
	trx.Status = "cancelled"

	return trx
}

var SetTransactionStatusInvalidAmount Transform = func(payment Payment, data interface{}) interface{} {
	trx := data.(Transaction)
	trx.Status = "invalid_amount"

	return trx
}
