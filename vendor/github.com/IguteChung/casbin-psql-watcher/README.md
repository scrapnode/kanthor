# Casbin Postgres Watcher
---
A WatcherEX implementation for [Casbin](https://github.com/casbin/casbin) based on [PostgreSQL](https://www.postgresql.org)

## Installation

  go get github.com/IguteChung/casbin-psql-watcher

## Example

```go
package main

import (
	"context"
	"fmt"

	psqlwatcher "github.com/IguteChung/casbin-psql-watcher"
	"github.com/casbin/casbin/v2"
)

func main() {
	// prepare the watcher with local psql server.
	// for demo purpose, enable NotifySelf to receive a local change callback.
	conn := "host=localhost user=postgres password=postgres dbname=postgres port=5432"
	w, _ := psqlwatcher.NewWatcherWithConnString(context.Background(), conn,
		psqlwatcher.Option{NotifySelf: true, Verbose: true})

	// prepare the enforcer.
	e, _ := casbin.NewEnforcer("../testdata/rbac_model.conf", "../testdata/rbac_policy.csv")

	// set the watcher for enforcer.
	_ = e.SetWatcher(w)

	// set the default callback to handle policy changes.
	_ = w.SetUpdateCallback(psqlwatcher.DefaultCallback(e))

	// update the policy and notify other enforcers.
	_ = e.SavePolicy()

	// wait for callback.
	fmt.Scanln()
}
```

## Dependency

- [Casbin](https://github.com/casbin/casbin)
- [pgx](https://github.com/jackc/pgx)
