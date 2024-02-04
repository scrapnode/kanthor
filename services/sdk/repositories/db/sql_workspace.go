package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	attributes := trace.WithAttributes(attribute.String("ws_id", id))
	subctx, span := ctx.Value(telemetry.CtxTracer).(trace.Tracer).Start(ctx, "repositories.db.workspace.get", attributes)
	defer func() {
		span.End()
	}()

	ws := &entities.Workspace{}

	transaction := database.SqlTxnFromContext(subctx, sql.client)
	tx := transaction.WithContext(subctx).Model(&ws).
		Where(fmt.Sprintf(`"%s"."id" = ?`, entities.TableWs), id).
		First(ws)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ws, nil
}
