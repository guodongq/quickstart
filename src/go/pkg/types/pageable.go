package types

type Pageable struct {
	Skip  int64 `json:"skip"`
	Limit int64 `json:"limit"`
}

func DefaultPageable() Pageable {
	return Pageable{
		Skip:  0,
		Limit: 100,
	}
}

type PageResult[T any] struct {
	Pageable     Pageable `json:"pageable"`
	TotalElement int64    `json:"total"`
	Content      []T      `json:"content"`
}
