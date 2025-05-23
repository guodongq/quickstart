package provider

import (
	"context"
	"fmt"
	"reflect"
	"time"

	logger "github.com/guodongq/quickstart/pkg/log"
)

func WaitForRunningProvider(p RunProvider, timeoutSeconds time.Duration) error {
	if p.IsRunning() {
		// No need to wait if provider is already running
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	name := Name(p)
	logger.WithField("timeout", timeoutSeconds).Debugf("Waiting for %s to run...", name)
	for {
		if p.IsRunning() {
			logger.Debugf("%s is running", name)
			return nil
		}

		if ctx.Err() != nil {
			return fmt.Errorf("time exceeded for %s to run", name)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func Name(provider Provider) string {
	return reflect.TypeOf(provider).Elem().String()
}
