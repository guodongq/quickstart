package eventbus_test

import (
	"context"
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"testing"

	. "github.com/guodongq/quickstart/pkg/util/eventbus"
	"github.com/stretchr/testify/require"
)

func TestEventBus(t *testing.T) {
	eventBus := New()

	var invoked bool
	zapLogger := logger.NewLogrusLogger()
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
