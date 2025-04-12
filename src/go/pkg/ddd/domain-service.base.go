package ddd

type DomainService interface {
}

type Factory[T any] interface {
	Create(...any) (T, error)
}
