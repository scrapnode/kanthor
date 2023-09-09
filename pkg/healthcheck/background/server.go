package background

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"os"
	"path"
	"time"
)

func NewServer(conf *healthcheck.Config, logger logging.Logger) healthcheck.Server {
	return &server{conf: conf, logger: logger, dest: path.Join(Dest, conf.Dest)}
}

type server struct {
	conf   *healthcheck.Config
	logger logging.Logger
	dest   string
}

func (server *server) Readiness(check func() error) error {
	if err := server.check(check); err != nil {
		return err
	}
	if err := server.write("readiness"); err != nil {
		return err
	}

	server.logger.Debug("ready")
	return nil
}

func (server *server) Liveness(check func() error) error {
	for {
		if err := server.check(check); err != nil {
			return err
		}

		if err := server.write("liveness"); err != nil {
			return err
		}

		server.logger.Debug("live", "timeout", server.conf.Timeout)
		time.Sleep(time.Millisecond * time.Duration(server.conf.Timeout))
	}
}

func (server *server) check(check func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(server.conf.Timeout))
	defer cancel()

	var err error
	for i := 0; i < server.conf.MaxTry; i++ {
		if err = check(); err == nil {
			return nil
		}
	}
	<-ctx.Done()
	return ctx.Err()
}

func (server *server) write(name string) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(time.Now().UTC().UnixMilli()))

	file := fmt.Sprintf("%s.%s", server.dest, name)
	return os.WriteFile(file, data, os.ModePerm)
}
