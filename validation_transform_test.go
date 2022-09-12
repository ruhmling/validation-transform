package validationAndTransform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_createValidationAlwaysTrue(t *testing.T) {
	var validation Validation
	var isValid bool

	validation = func(payment Payment, x interface{}) bool {
		return true
	}

	isValid = validation(Payment{}, 10)

	assert.Equal(t, true, isValid)
}

func Test_createTransformNothing(t *testing.T) {
	var transform Transform = func(payment Payment, data interface{}) interface{} {
		return data
	}

	dataTransformed := transform(Payment{}, 10)

	assert.Equal(t, 10, dataTransformed)
}

func Test_givenPaymentWithAmount10AndTotalAmount15ThenReturnInvalidAmount(t *testing.T) {
	payment := Payment{Amount: 10}
	totalAmount := float64(15)

	assert.Equal(t, false, IsAmountValid(payment, totalAmount))
}

func Test_givenPaymentWithAmount15AndTotalAmount15ThenReturnValidAmount(t *testing.T) {
	payment := Payment{Amount: 15}
	totalAmount := float64(15)

	assert.Equal(t, true, IsAmountValid(payment, totalAmount))
}

func Test_givenPaymentWithAmount10AndTransactionWithAmount15ThenReturnInvalidAmount(t *testing.T) {
	payment := Payment{Amount: 10}
	trx := Transaction{Amount: 15}

	assert.Equal(t, false, IsTransactionAmountValid(payment, trx))
}

func Test_givenPaymentWithAmount15AndTransactionWithAmount15ThenReturnValidAmount(t *testing.T) {
	payment := Payment{Amount: 15}
	trx := Transaction{Amount: 15}

	assert.Equal(t, true, IsTransactionAmountValid(payment, trx))
}

func Test_givenPaymentWithAmount10AndTransactionWithAmount9And50ThenSubsidizeAmount(t *testing.T) {
	payment := Payment{Amount: 10}
	trx := Transaction{Amount: 9.5}

	if IsAmountSubsidisable(payment, trx) {
		trx = SubsidizeTransaction(payment, trx).(Transaction)
	}

	assert.Equal(t, 10.0, trx.Amount)
}

func Test_givenPaymentWithAmount11AndTransactionWithAmount10And50ThenSubsidizeAmount(t *testing.T) {
	payment := Payment{Amount: 11}
	trx := Transaction{Amount: 10.5}

	if IsAmountSubsidisable(payment, trx) {
		trx = SubsidizeTransaction(payment, trx).(Transaction)
	}

	assert.Equal(t, 11.0, trx.Amount)
}

func Test_givenPaymentWithStatusCancelledAndTransactionToProcessThenReturnCancelledTransaction(t *testing.T) {
	payment := Payment{Status: "cancelled"}
	trx := Transaction{Status: "to_process"}

	if IsPaymentCancelled(payment, trx) {
		trx = SetTransactionStatusCancelled(payment, trx).(Transaction)
	}

	assert.Equal(t, "cancelled", trx.Status)
}

func Test_givenPaymentWithAmount11AndTransactionWithAmount9ThenReturnStatusInvalidAmount(t *testing.T) {
	payment := Payment{Amount: 11}
	trx := Transaction{Amount: 9}

	if IsTransactionAmountInsufficient(payment, trx) {
		trx = SetTransactionStatusInvalidAmount(payment, trx).(Transaction)
	}

	assert.Equal(t, "invalid_amount", trx.Status)
}

func Test_createValidationAndTransformHandler(t *testing.T) {
	assert.NotNil(t, validationAndTransform{})
}

func Test_executeInsufficientAmountValidationThenReturnInvalidAmountTransaction(t *testing.T) {
	validationAndTransform := validationAndTransform{
		Validation: IsTransactionAmountInsufficient,
		Transform:  SetTransactionStatusInvalidAmount,
	}

	trx := validationAndTransform.Execute(Payment{Amount: 10}, Transaction{Amount: 9}).(Transaction)

	assert.Equal(t, "invalid_amount", trx.Status)
}

func Test_executeCancelledPaymentValidationThenReturnCancelledTransaction(t *testing.T) {
	validationAndTransform := validationAndTransform{
		Validation: IsPaymentCancelled,
		Transform:  SetTransactionStatusCancelled,
	}

	trx := validationAndTransform.Execute(Payment{Status: "cancelled"}, Transaction{}).(Transaction)

	assert.Equal(t, "cancelled", trx.Status)
}

func Test_executeValidationsAndTransformHandler(t *testing.T) {
	validationAndTransform := NewValidationAndTransformBuilder(IsPaymentCancelled, SetTransactionStatusCancelled, true).
		AddNext(IsAmountSubsidisable, SubsidizeTransaction, false).
		AddNext(IsTransactionAmountInsufficient, SetTransactionStatusInvalidAmount, true).
		Build()

	trx := validationAndTransform.Execute(Payment{Amount: 10, Status: "pending"}, Transaction{Amount: 10, Status: "to_process"}).(Transaction)

	assert.Equal(t, "to_process", trx.Status)
	assert.Equal(t, 10.0, trx.Amount)

	trx = validationAndTransform.Execute(Payment{Amount: 10, Status: "pending"}, Transaction{Amount: 9.5, Status: "to_process"}).(Transaction)

	assert.Equal(t, "amount_subsidized", trx.StatusDetail)
	assert.Equal(t, "to_process", trx.Status)
	assert.Equal(t, 10.0, trx.Amount)

	trx = validationAndTransform.Execute(Payment{Amount: 10, Status: "pending"}, Transaction{Amount: 8, Status: "to_process"}).(Transaction)

	assert.Equal(t, "", trx.StatusDetail)
	assert.Equal(t, "invalid_amount", trx.Status)
	assert.Equal(t, 8.0, trx.Amount)

	trx = validationAndTransform.Execute(Payment{Amount: 10, Status: "cancelled"}, Transaction{Amount: 10, Status: "to_process"}).(Transaction)

	assert.Equal(t, "", trx.StatusDetail)
	assert.Equal(t, "cancelled", trx.Status)
	assert.Equal(t, 10.0, trx.Amount)

	trx = validationAndTransform.Execute(Payment{Amount: 10, Status: "cancelled"}, Transaction{Amount: 9.5, Status: "to_process"}).(Transaction)

	assert.Equal(t, "", trx.StatusDetail)
	assert.Equal(t, "cancelled", trx.Status)
	assert.Equal(t, 9.5, trx.Amount)
}
