package background

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
)

var Readiness = "readiness"
var Liveness = "liveness"

func NewServer(conf *healthcheck.Config, logger logging.Logger) healthcheck.Server {
	return &server{
		conf:       conf,
		logger:     logger,
		dest:       path.Join(Dest, conf.Dest),
		terminated: make(chan int64, 1),
	}
}

type server struct {
	conf   *healthcheck.Config
	logger logging.Logger
	dest   string

	terminated chan int64
}

func (server *server) Connect(ctx context.Context) error {
	server.logger.Info("HEALTHCHECK.BACKGROUND.SERVER.CONNECTED")
	return nil
}

func (server *server) Disconnect(ctx context.Context) error {
	server.terminated <- time.Now().UTC().UnixMilli()
	server.logger.Info("HEALTHCHECK.BACKGROUND.SERVER.DISCONNECTED")
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

func (server *server) Liveness(check func() error) error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(server.conf.Liveness.Timeout))
	defer ticker.Stop()

	for {
		select {
		case <-server.terminated:
			return nil
		case <-ticker.C:
			if err := server.check(Liveness, &server.conf.Liveness, check); err != nil {
				return err
			}
			if err := server.write(Liveness); err != nil {
				return err
			}
			server.logger.Debug("live")
		}
	}
}

func (server *server) check(name string, conf *healthcheck.CheckConfig, check func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.Timeout))
	defer cancel()

	errc := make(chan error, 1)
	go func() {
		var returning error
		for i := 0; i < conf.MaxTry; i++ {
			err := check()
			if err == nil {
				errc <- nil
				return
			}

			returning = errors.Join(returning, err)
		}

		errc <- fmt.Errorf("HEALTHCHECK.BACKGROUND.SERVER.ERROR: %v", returning)
	}()

	select {
	case <-server.terminated:
		return nil
	case err := <-errc:
		return err
	case <-ctx.Done():
		return fmt.Errorf("HEALTHCHECK.BACKGROUND.SERVER.ERROR: %v | timeout:%d max_try:%d", ctx.Err(), conf.Timeout, conf.MaxTry)
	}
}

func (server *server) write(name string) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(time.Now().UTC().UnixMilli()))

	file := fmt.Sprintf("%s.%s", server.dest, name)
	return os.WriteFile(file, data, os.ModePerm)
}
