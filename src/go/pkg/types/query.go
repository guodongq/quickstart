package types

type QueryOption func(*QueryOptions)

type QueryOptions struct {
	Pageable
	Sortable
	Filter
}

func NewQueryOptions(opts ...QueryOption) *QueryOptions {
	q := &QueryOptions{
		Pageable: DefaultPageable(),
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

func WithPageable(skip, limit int64) QueryOption {
	return func(q *QueryOptions) {
		q.Skip = skip
		q.Limit = limit
	}
}

func WithSort(sort ...string) QueryOption {
	return func(q *QueryOptions) {
		q.Sort = append(q.Sort, sort...)
	}
}

func WithFilter(filter Filter) QueryOption {
	return func(q *QueryOptions) {
		q.Filter = filter
	}
}
