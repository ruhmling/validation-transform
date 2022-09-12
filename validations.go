package validationAndTransform

var IsAmountValid Validation = func(payment Payment, data interface{}) bool {
	return payment.Amount == data.(float64)
}

var IsTransactionAmountValid Validation = func(payment Payment, data interface{}) bool {
	trx := data.(Transaction)

	return payment.Amount == trx.Amount
}

var IsTransactionAmountInsufficient Validation = func(payment Payment, data interface{}) bool {
	trx := data.(Transaction)

	return payment.Amount > trx.Amount
}

var IsAmountSubsidisable Validation = func(payment Payment, data interface{}) bool {
	trx := data.(Transaction)

	return trx.Amount < payment.Amount && payment.Amount-trx.Amount < 1
}

var IsPaymentCancelled Validation = func(payment Payment, data interface{}) bool {
	return payment.Status == "cancelled"
}
