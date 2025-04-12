package app

import (
	"os"
	"strings"
)

type AppOptions struct {
	Name     string
	BasePath string
}

func getDefaultAppOptions() AppOptions {
	paths := strings.Split(os.Args[0], "/")
	appName := paths[len(paths)-1]

	return AppOptions{
		Name:     appName,
		BasePath: "/",
	}
}

func WithAppOptionsName(name string) func(*AppOptions) {
	return func(o *AppOptions) {
		o.Name = name
	}
}

func WithAppOptionsBasePath(basePath string) func(*AppOptions) {
	return func(o *AppOptions) {
		o.BasePath = basePath
	}
}
