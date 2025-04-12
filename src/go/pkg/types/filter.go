package types

type F map[string]any

type Filter interface {
	ToFilter() F
	Or(...Filter) Filter
	And(...Filter) Filter
}

type DefaultMongoFilter struct {
	F
}

func (m *DefaultMongoFilter) ToFilter() F {
	return m.F
}

func (m *DefaultMongoFilter) Or(alternativeFilters ...Filter) Filter {
	if m.F == nil {
		m.F = F{}
	}

	var alternatives []F
	for _, alternativeFilter := range alternativeFilters {
		alternatives = append(alternatives, alternativeFilter.ToFilter())
	}

	m.F["$or"] = alternatives
	return m
}

func (m *DefaultMongoFilter) And(alternativeFilters ...Filter) Filter {
	if m.F == nil {
		m.F = F{}
	}

	var alternatives []F
	for _, alternativeFilter := range alternativeFilters {
		alternatives = append(alternatives, alternativeFilter.ToFilter())
	}

	m.F["$and"] = alternatives
	return m
}
