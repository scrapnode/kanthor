package dataplane

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

type worksapce struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
