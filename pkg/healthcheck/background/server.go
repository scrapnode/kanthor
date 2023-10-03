package background

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
)

func NewServer(conf *healthcheck.Config, logger logging.Logger) healthcheck.Server {
	return &server{conf: conf, logger: logger, dest: path.Join(Dest, conf.Dest)}
}

var (
	StatusStopped = -1
)

type server struct {
	conf   *healthcheck.Config
	logger logging.Logger
	dest   string
	status int
}

func (server *server) Connect(ctx context.Context) error {
	return nil
}

func (server *server) Disconnect(ctx context.Context) error {
	server.status = StatusStopped
	return nil
}

func (server *server) Readiness(check func() error) error {
	if err := server.check(&server.conf.Readiness, check); err != nil {
		return err
	}
	if err := server.write("readiness"); err != nil {
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
		if server.status == StatusStopped {
			return
		}

		if err = server.check(&server.conf.Liveness, check); err != nil {
			return
		}

		if err = server.write("liveness"); err != nil {
			return
		}

		server.logger.Debugw("live", "timeout", server.conf.Liveness.Timeout)
		time.Sleep(time.Millisecond * time.Duration(server.conf.Liveness.Timeout))
	}
}

func (server *server) check(conf *healthcheck.CheckConfig, check func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.Timeout))
	defer cancel()

	var err error
	for i := 0; i < conf.MaxTry; i++ {
		if err = check(); err == nil {
			return nil
		}
	}
	<-ctx.Done()
	return fmt.Errorf("HEALTHCHECK.BACKGROUND.ERROR: %v | timeout:%d max_try:%d", ctx.Err(), conf.Timeout, conf.MaxTry)
}

func (server *server) write(name string) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(time.Now().UTC().UnixMilli()))

	file := fmt.Sprintf("%s.%s", server.dest, name)
	return os.WriteFile(file, data, os.ModePerm)
}
