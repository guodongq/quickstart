package bus_test

import (
	"context"
	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider/logger/zap"
	"testing"

	. "github.com/guodongq/quickstart/pkg/bus"
	"github.com/stretchr/testify/require"
)

type ProjectCreatedEvent struct {
	ProjectID   int32
	ProjectName string
}

func TestEventBus(t *testing.T) {
	eventBus := New()

	var invoked bool
	zapLogger := zap.New()
	zapLogger.Init()

	eventBus.Subscribe(func(_ context.Context, event *ProjectCreatedEvent) error {
		logger.Info("Event1: coming request", event)
		invoked = true
		return nil
	})

	eventBus.Subscribe(func(_ context.Context, event *ProjectCreatedEvent) error {
		logger.Info("Event2: coming request", event)
		invoked = true
		return nil
	})

	err := eventBus.Publish(context.Background(), &ProjectCreatedEvent{ProjectID: 1, ProjectName: "abc"})
	require.NoError(t, err, "unable to publish event")
	require.True(t, invoked)
}
