package main

import (
	cli "github.com/guodongq/quickstart/cmd/cli-quickstart/commands"
	gogrpc "github.com/guodongq/quickstart/cmd/go-grpc-quickstart/commands"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const binaryNameEnv = "QUICKSTART_BINARY_NAME"

func main() {
	var command *cobra.Command

	binaryName := filepath.Base(os.Args[0])
	if val := os.Getenv(binaryNameEnv); val != "" {
		binaryName = val
	}

	switch binaryName {
	case "go-grpc-quickstart":
		command = gogrpc.NewCommand()
	default:
		command = cli.NewCommand()
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
