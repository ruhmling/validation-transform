package validationAndTransform

func NewValidationAndTransformBuilder(validation Validation, transform Transform, isBreakPoint bool) *builder {
	builder := new(builder)
	builder.list = []Handler{
		newValidationAndTransform(validation, transform, isBreakPoint),
	}

	return builder
}

type builder struct {
	list []Handler
}

func (v *builder) AddNext(validation Validation, transform Transform, isBreakPoint bool) *builder {
	v.list = append(v.list, newValidationAndTransform(validation, transform, isBreakPoint))

	return v
}

func (v *builder) Build() Handler {
	first := v.list[0]
	next := first
	size := len(v.list)

	for i := 0; i < size; i++ {
		if i+1 == size {
			break
		}

		next.AddNext(v.list[i+1])
		next = v.list[i+1]
	}

	return first
}

func newValidationAndTransform(validation Validation, transform Transform, isBreakPoint bool) Handler {
	valAndTrans := validationAndTransform{
		Validation:   validation,
		Transform:    transform,
		IsBreakPoint: isBreakPoint,
	}

	return &valAndTrans
}

type validationAndTransform struct {
	Validation   Validation
	Transform    Transform
	IsBreakPoint bool
	Next         Handler
}

func (v *validationAndTransform) AddNext(handler Handler) {
	v.Next = handler
}

func (v *validationAndTransform) Execute(payment Payment, data interface{}) interface{} {
	transaction := data.(Transaction)
	isTransformed := false

	if v.Validation(payment, transaction) {
		transaction = v.Transform(payment, transaction).(Transaction)
		isTransformed = true
	}

	if (v.IsBreakPoint && isTransformed) || v.Next == nil {
		return transaction
	}

	return v.Next.Execute(payment, transaction)
}
