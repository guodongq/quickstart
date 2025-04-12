package ddd

import "reflect"

type ValueObject interface {
	Equals(ValueObject) bool
}

type BaseValueObject struct{}

func (vo *BaseValueObject) Equals(other ValueObject) bool {
	// Need specific implementation for each value object
	return reflect.DeepEqual(vo, other)
}
