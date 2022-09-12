package validationAndTransform

type Validation func(payment Payment, data interface{}) bool
type Transform func(payment Payment, data interface{}) interface{}

type Handler interface {
	AddNext(handler Handler)
	Execute(payment Payment, data interface{}) interface{}
}
