package migrate

import (
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
)

func Source(task *config.MigrationTask, logger logging.Logger) (migration.Source, error) {
	if task.Name == "database" {
		return database.New(&database.Config{Uri: task.Uri}, logger), nil
	}
	if task.Name == "datastore" {
		return datastore.New(&datastore.Config{Uri: task.Uri}, logger), nil
	}

	return nil, fmt.Errorf("migrate: unsupport task [%s]", task.Name)
}
