package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/status"
)

func (uc *endeavor) Evaluate(ctx context.Context, attempts []entities.Attempt) (*entities.AttemptStrive, error) {
	returning := &entities.AttemptStrive{Attemptable: []entities.Attempt{}, Ignore: []string{}}

	for _, attempt := range attempts {
		if attempt.Status == status.ErrIgnore {
			continue
		}
		if attempt.Complete() {
			continue
		}

		if attempt.Status == status.ErrUnknown {
			returning.Attemptable = append(returning.Attemptable, attempt)
			continue
		}

		if attempt.Status == status.None {
			returning.Attemptable = append(returning.Attemptable, attempt)
			continue
		}

		if status.Is5xx(attempt.Status) {
			returning.Attemptable = append(returning.Attemptable, attempt)
			continue
		}

		returning.Ignore = append(returning.Ignore, attempt.ReqId)
		uc.logger.Warnw("ignore attempt", "req_id", attempt.ReqId, "status", attempt.Status)
	}

	uc.logger.Debugw("evaluate records", "attempt_count", len(returning.Attemptable), "Ignore_count", len(returning.Ignore))
	return returning, nil
}
