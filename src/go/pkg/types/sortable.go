package types

import (
	"github.com/guodongq/quickstart/pkg/util/log"
	"strings"
)

var symbolDirectionMapping = map[string]Direction{
	"-": DESC,
	"+": ASC,
}

type Direction string

const (
	ASC  Direction = "asc"
	DESC Direction = "desc"
)

type Sort struct {
	OrderBy   string
	Direction Direction
}

type Sortable struct {
	Sort []string
}

func DefaultSortable(sortableFuncs ...func(sortable *Sortable)) Sortable {
	var d Sortable
	for _, sortableFunc := range sortableFuncs {
		sortableFunc(&d)
	}
	return d
}

func (s *Sortable) ToSort() any {
	return s.ComputedSort()
}

func (s *Sortable) ComputedSort() []Sort {
	var result []Sort
	for _, item := range s.Sort {
		sort, err := parseSort(item)
		if err != nil {
			log.Infof("Failed to parse sort item '%s': %v\n", item, err)
			continue
		}
		if sort.OrderBy != "" {
			result = append(result, sort)
		}
	}
	return result
}

func parseSort(item string) (Sort, error) {
	item = strings.TrimSpace(item)
	if len(item) == 0 {
		return Sort{}, nil
	}

	var direction Direction
	var orderBy string

	if symbol, exists := symbolDirectionMapping[string(item[0])]; exists {
		direction = symbol
		orderBy = item[1:]
	} else {
		direction = ASC
		orderBy = item
	}

	return Sort{OrderBy: orderBy, Direction: direction}, nil
}
