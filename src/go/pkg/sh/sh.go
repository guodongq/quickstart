package sh

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Verbose enables verbose output
var Verbose bool

// RunCommand executes a shell command.
func RunCommand(name string, arg ...string) error {
	if Verbose {
		cmdText := name + " " + strings.Join(arg, " ")
		fmt.Fprintln(os.Stderr, " + ", cmdText)
	}
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SplitParameters splits shell command parameters, taking quoting in account.
func SplitParameters(s string) []string {
	r := regexp.MustCompile(`'[^']*'|[^ ]+`)
	return r.FindAllString(s, -1)
}
