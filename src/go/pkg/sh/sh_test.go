package sh

import (
	"strings"
	"testing"

	logger "github.com/guodongq/quickstart/pkg/log"
)

func TestSplitParameters(t *testing.T) {
	in := `-a -tags 'netgo static_build'`
	expect := []string{"-a", "-tags", `'netgo static_build'`}
	got := SplitParameters(in)
	for i, g := range got {
		if expect[i] != g {
			t.Error("expected", expect[i], "got", g, "full output: ", strings.Join(got, "#"))
		}
	}
}

func TestRunCommand(t *testing.T) {
	Verbose = true
	err := RunCommand("echo", "hello")
	logger.Error("error")
	if err != nil {
		t.Error("expected nil got", err)
	}
}
