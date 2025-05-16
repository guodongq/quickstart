package main

import (
	"os"
	"path/filepath"

	academy "github.com/guodongq/quickstart/cmd/academy/commands"
	cli "github.com/guodongq/quickstart/cmd/cli/commands"
	gogrpc "github.com/guodongq/quickstart/cmd/go-grpc/commands"
	"github.com/spf13/cobra"
)

const binaryNameEnv = "QUICKSTART_BINARY_NAME"

func init() {
	_ = os.Setenv(binaryNameEnv, "academy")
}

func main() {
	var command *cobra.Command

	binaryName := filepath.Base(os.Args[0])
	if val := os.Getenv(binaryNameEnv); val != "" {
		binaryName = val
	}

	switch binaryName {
	case "go-grpc":
		command = gogrpc.NewCommand()
	case "academy":
		command = academy.NewCommand()
	default:
		command = cli.NewCommand()
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
