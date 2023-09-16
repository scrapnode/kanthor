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
// set KANTHOR_DEBUGGING_SERVER=true to enable it
func NewServer() Server {
	enable := strings.EqualFold(os.Getenv("KANTHOR_DEBUGGING_SERVER"), "true")
	return &server{enable: enable}
}

type server struct {
	enable   bool
	instance *http.Server
}

func (server *server) Start(ctx context.Context) error {
	if !server.enable {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server.instance = &http.Server{Addr: fmt.Sprintf(":%d", HTTP_PORT), Handler: mux}
	return nil
}

func (server *server) Stop(ctx context.Context) error {
	if !server.enable {
		return nil
	}

	return server.instance.Shutdown(ctx)
}

func (server *server) Run(ctx context.Context) error {
	if !server.enable {
		return nil
	}

	return server.instance.ListenAndServe()
}
