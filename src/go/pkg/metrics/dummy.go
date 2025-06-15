package metrics

import "github.com/guodongq/quickstart/pkg/decorator"

type NoOp struct{}

func New() decorator.MetricsClient {
	return NoOp{}
}

func (n NoOp) Inc(_ string, _ int) {
	// todo - add some implementation!
}
