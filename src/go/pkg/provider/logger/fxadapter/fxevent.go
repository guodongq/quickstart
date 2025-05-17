package fxadapter

import (
	"fmt"
	"github.com/guodongq/quickstart/pkg/log"
	"go.uber.org/fx/fxevent"
	"strings"
)

type FxEventLogger struct {
	log.Logger
}

func NewFxEventLogger(logger log.Logger) *FxEventLogger {
	return &FxEventLogger{
		Logger: logger,
	}
}

func (l *FxEventLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Infof("HOOK OnStart    %s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Infof("HOOK OnStart    %s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			l.Infof("HOOK OnStart    %s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.OnStopExecuting:
		l.Infof("HOOK OnStop    %s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Infof("HOOK OnStop    %s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			l.Infof("HOOK OnStop    %s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Infof("ERROR  Failed to supply %v: %+v", e.TypeName, e.Err)
		} else if e.ModuleName != "" {
			l.Infof("SUPPLY  %v from module %q", e.TypeName, e.ModuleName)
		} else {
			l.Infof("SUPPLY  %v", e.TypeName)
		}
	case *fxevent.Provided:
		var privateStr string
		if e.Private {
			privateStr = " (PRIVATE)"
		}
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.Infof("PROVIDE%v  %v <= %v from module %q", privateStr, rtype, e.ConstructorName, e.ModuleName)
			} else {
				l.Infof("PROVIDE%v  %v <= %v", privateStr, rtype, e.ConstructorName)
			}
		}
		if e.Err != nil {
			l.Infof("Error after options were applied: %+v", e.Err)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.Infof("REPLACE  %v from module %q", rtype, e.ModuleName)
			} else {
				l.Infof("REPLACE  %v", rtype)
			}
		}
		if e.Err != nil {
			l.Infof("ERROR  Failed to replace: %+v", e.Err)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.Infof("DECORATE  %v <= %v from module %q", rtype, e.DecoratorName, e.ModuleName)
			} else {
				l.Infof("DECORATE  %v <= %v", rtype, e.DecoratorName)
			}
		}
		if e.Err != nil {
			l.Infof("Error after options were applied: %+v", e.Err)
		}
	case *fxevent.Run:
		var moduleStr string
		if e.ModuleName != "" {
			moduleStr = fmt.Sprintf(" from module %q", e.ModuleName)
		}
		l.Infof("RUN  %v: %v in %s%v", e.Kind, e.Name, e.Runtime, moduleStr)
		if e.Err != nil {
			l.Infof("Error returned: %+v", e.Err)
		}

	case *fxevent.Invoking:
		if e.ModuleName != "" {
			l.Infof("INVOKE    %s from module %q", e.FunctionName, e.ModuleName)
		} else {
			l.Infof("INVOKE    %s", e.FunctionName)
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Infof("ERROR    fx.Invoke(%v) called from:\n%+vFailed: %+v", e.FunctionName, e.Trace, e.Err)
		}
	case *fxevent.Stopping:
		l.Infof("%v", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Infof("ERROR    Failed to stop cleanly: %+v", e.Err)
		}
	case *fxevent.RollingBack:
		l.Infof("ERROR    Start failed, rolling back: %+v", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Infof("ERROR    Couldn't roll back cleanly: %+v", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Infof("ERROR    Failed to start: %+v", e.Err)
		} else {
			l.Infof("RUNNING")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Infof("ERROR    Failed to initialize custom logger: %+v", e.Err)
		} else {
			l.Infof("LOGGER  Initialized custom logger from %v", e.ConstructorName)
		}
	}
}
