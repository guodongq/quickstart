package main

import (
	"github.com/guodongq/quickstart/pkg/stack"
)

func main() {
	st := stack.New()
	defer st.MustClose()

	// todo: add your code here

	st.MustRun()
}
