package services

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Service interface {
	patterns.Runnable
}
