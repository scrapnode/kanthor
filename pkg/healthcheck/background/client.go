package background

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/scrapnode/kanthor/pkg/healthcheck"
)

func NewClient(conf *healthcheck.Config) healthcheck.Client {
	return &client{conf: conf, dest: path.Join(Dest, conf.Dest)}
}

type client struct {
	conf *healthcheck.Config
	dest string
}

func (client *client) Readiness() error {
	_, err := client.read("readiness")
	return err
}

func (client *client) Liveness() error {
	diff, err := client.read("liveness")
	if err != nil {
		return err
	}

	delta := int64(client.conf.Liveness.Timeout * client.conf.Liveness.MaxTry)
	if diff > delta {
		return fmt.Errorf("timeout (diff:%d delta:%d)", diff, delta)
	}

	return nil
}

func (client *client) read(name string) (int64, error) {
	file := fmt.Sprintf("%s.%s", client.dest, name)
	data, err := os.ReadFile(file)
	if err != nil {
		return 0, err
	}

	prev := int64(binary.BigEndian.Uint64(data))
	now := time.Now().UnixMilli()
	return now - prev, nil
}
