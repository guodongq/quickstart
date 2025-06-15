package chi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/middleware/requestid"
	"github.com/guodongq/quickstart/pkg/provider"
)

type ChiEngine struct {
	provider.AbstractRunProvider

	srv *http.Server

	options ChiOptions
}

func New(optionsFuncs ...func(options *ChiOptions)) *ChiEngine {
	defaultOptions := getDefaultChiEngineOptions()
	options := &defaultOptions

	for _, optionsFunc := range optionsFuncs {
		optionsFunc(options)
	}

	return &ChiEngine{
		options: *options,
	}
}

func configureMiddlewares(router chi.Router, options *ChiOptions) http.Handler {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(logger.NewStructuredLogger(logger.DefaultLogger()))
	router.Use(Logger(true))
	router.Use(middleware.Recoverer)
	//router.Mount("/debug", middleware.Profiler())
	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)

	for _, v := range options.Middlewares {
		router.Use(v)
	}

	var handler http.Handler = router
	if options.HandlerFromMux != nil {
		handler = options.HandlerFromMux(router)
	}
	return handler
}

func (p *ChiEngine) Run() error {
	addr := fmt.Sprintf(":%d", p.options.Port)

	logEntry := logger.WithFields(logger.Fields{
		"addr": addr,
	})

	handler := configureMiddlewares(chi.NewRouter(), &p.options)

	p.srv = &http.Server{
		Handler: handler,
		Addr:    addr,
	}

	p.SetRunning(true)

	logEntry.Info("Chi Engine Provider Launched")
	if err := p.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logEntry.WithError(err).Error("Chi Engine Provider launch failed")
		return err
	}

	return nil
}

func Logger(useColor bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			buf := &bytes.Buffer{}
			reqID := requestid.Get(r)
			if reqID != "" {
				cW(buf, useColor, nYellow, "[%s] ", reqID)
			}
			cW(buf, useColor, nCyan, "")
			cW(buf, useColor, bMagenta, "%s ", r.Method)

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			cW(buf, useColor, nCyan, "%s://%s%s %s ", scheme, r.Host, r.RequestURI, r.Proto)

			buf.WriteString("from ")
			buf.WriteString(r.RemoteAddr)
			buf.WriteString(" - ")

			logEntry := logger.WithFields(logger.Fields{
				"request": buf.String(),
			})
			logEntry.Info("Request started")

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				logEntry.WithField("response", fmt.Sprintf("%d %d %s %v",
					ww.Status(),
					ww.BytesWritten(),
					ww.Header(),
					time.Since(t1)),
				).Info("Request completed")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// colorWrite
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if IsTTY && useColor {
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if IsTTY && useColor {
		w.Write(reset)
	}
}

var IsTTY bool

func init() {
	// This is sort of cheating: if stdout is a character device, we assume
	// that means it's a TTY. Unfortunately, there are many non-TTY
	// character devices, but fortunately stdout is rarely set to any of
	// them.
	//
	// We could solve this properly by pulling in a dependency on
	// code.google.com/p/go.crypto/ssh/terminal, for instance, but as a
	// heuristic for whether to print in color or in black-and-white, I'd
	// really rather not.
	fi, err := os.Stdout.Stat()
	if err == nil {
		m := os.ModeDevice | os.ModeCharDevice
		IsTTY = fi.Mode()&m == m
	}
}

var (
	// Normal colors
	nBlack   = []byte{'\033', '[', '3', '0', 'm'}
	nRed     = []byte{'\033', '[', '3', '1', 'm'}
	nGreen   = []byte{'\033', '[', '3', '2', 'm'}
	nYellow  = []byte{'\033', '[', '3', '3', 'm'}
	nBlue    = []byte{'\033', '[', '3', '4', 'm'}
	nMagenta = []byte{'\033', '[', '3', '5', 'm'}
	nCyan    = []byte{'\033', '[', '3', '6', 'm'}
	nWhite   = []byte{'\033', '[', '3', '7', 'm'}
	// Bright colors
	bBlack   = []byte{'\033', '[', '3', '0', ';', '1', 'm'}
	bRed     = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen   = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	bYellow  = []byte{'\033', '[', '3', '3', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	bWhite   = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)
