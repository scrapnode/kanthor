package background

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
)

var Readiness = "readiness"
var Liveness = "liveness"

func NewServer(conf *healthcheck.Config, logger logging.Logger) healthcheck.Server {
	return &server{conf: conf, logger: logger, dest: path.Join(Dest, conf.Dest)}
}

type server struct {
	conf       *healthcheck.Config
	logger     logging.Logger
	dest       string
	terminated bool
}

func (server *server) Connect(ctx context.Context) error {
	return nil
}

func (server *server) Disconnect(ctx context.Context) error {
	server.terminated = true
	return nil
}

func (server *server) Readiness(check func() error) error {
	if err := server.check(Readiness, &server.conf.Readiness, check); err != nil {
		return err
	}
	if err := server.write(Readiness); err != nil {
		return err
	}

	server.logger.Debug("ready")
	return nil
}

func (server *server) Liveness(check func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			server.logger.Errorf("healthcheck.background.liveness.recover: %v", r)
			server.logger.Errorf("healthcheck.background.liveness.recover.trace: %s", debug.Stack())

			if rerr, ok := r.(error); ok {
				err = rerr
				return
			}
		}
	}()

	for {
		if server.terminated {
			return
		}

		if err = server.check(Liveness, &server.conf.Liveness, check); err != nil {
			return
		}

		if err = server.write(Liveness); err != nil {
			return
		}

		server.logger.Debugw("live", "timeout", server.conf.Liveness.Timeout)
		time.Sleep(time.Millisecond * time.Duration(server.conf.Liveness.Timeout))
	}
}

func (server *server) check(name string, conf *healthcheck.CheckConfig, check func() error) error {
	for i := 0; i < conf.MaxTry; i++ {
		time.Sleep(time.Millisecond * time.Duration(conf.Timeout/conf.MaxTry))

		if server.terminated {
			return nil
		}
		err := check()
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("HEALTHCHECK.BACKGROUND.ERROR| timeout:%d max_try:%d", conf.Timeout, conf.MaxTry)
}

func (server *server) write(name string) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(time.Now().UTC().UnixMilli()))

	file := fmt.Sprintf("%s.%s", server.dest, name)
	return os.WriteFile(file, data, os.ModePerm)
}
