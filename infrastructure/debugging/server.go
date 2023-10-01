package debugging

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

var HTTP_PORT = 6060

type Server interface {
	patterns.Runnable
}

// NewServer of debugging can only enable via environment variable
// set KANTHOR_DEBUGGING_SERVER_ENABLE=true to enable it
func NewServer() Server {
	enable := strings.EqualFold(os.Getenv("KANTHOR_DEBUGGING_SERVER_ENABLE"), "true")
	return &server{enable: enable}
}

type server struct {
	enable   bool
	instance *http.Server
}

func (server *server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		text := "disable"
		if server.enable {
			text = "enable"
		}
		_, _ = w.Write([]byte(text))
	})

	if server.enable {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	addr := fmt.Sprintf(":%d", HTTP_PORT)
	if port := os.Getenv("KANTHOR_DEBUGGING_SERVER_PORT"); port != "" {
		addr = fmt.Sprintf(":%s", port)
	}
	server.instance = &http.Server{Addr: addr, Handler: mux}
	return nil
}

func (server *server) Stop(ctx context.Context) error {
	return server.instance.Shutdown(ctx)
}

func (server *server) Run(ctx context.Context) error {
	return server.instance.ListenAndServe()
}
