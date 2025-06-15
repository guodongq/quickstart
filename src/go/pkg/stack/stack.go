package stack

import (
	"os"
	"os/signal"
	"sync"

	logger "github.com/guodongq/quickstart/pkg/log"

	p "github.com/guodongq/quickstart/pkg/provider"
)

var (
	runOnce   sync.Once
	closeOnce sync.Once
)

type Stack struct {
	providers []p.Provider
}

func New() *Stack {
	return &Stack{
		providers: make([]p.Provider, 0),
	}
}

func (s *Stack) MustInit(provider p.Provider) {
	name := p.Name(provider)
	logger.Debugf("%s initializing...", name)

	if err := provider.Init(); err != nil {
		logger.WithError(err).Panicf("Error during %s initialization", name)
	}

	s.providers = append(s.providers, provider)
	logger.Infof("%s initialized", name)
}

func (s *Stack) MustClose() {
	closeOnce.Do(func() {
		for i := len(s.providers) - 1; i >= 0; i-- {
			name := p.Name(s.providers[i])
			logger.Debugf("%s closing...", name)

			if err := s.providers[i].Close(); err != nil {
				logger.WithError(err).Panicf("%s failed to close", name)
			}

			logger.Infof("%s closed", name)
		}
	})
}

func (s *Stack) MustRun() {
	runOnce.Do(func() {
		for _, provider := range s.providers {
			if runProvider, ok := provider.(p.RunProvider); ok {
				go s.launch(runProvider)
			}
		}
		s.handleInterrupt()
	})
}

func (s *Stack) handleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		s.MustClose()
		close(cleanupDone)
	}()
	<-cleanupDone
}

func (s *Stack) launch(provider p.RunProvider) {
	name := p.Name(provider)
	logger.Debugf("%s launching...", name)

	if err := provider.Run(); err != nil {
		logger.WithError(err).Panicf("%s failed to run", name)
	}

	logger.Debugf("%s launched", name)
}
