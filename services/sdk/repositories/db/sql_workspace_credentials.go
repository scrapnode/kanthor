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

type SqlWorkspaceCredentials struct {
	client *gorm.DB
}

func (sql *SqlWorkspaceCredentials) Get(ctx context.Context, id string) (*entities.WorkspaceCredentials, error) {
	attributes := trace.WithAttributes(attribute.String("wsc.id", id))
	subctx, span := ctx.Value(telemetry.CtxTracer).(trace.Tracer).Start(ctx, "repositories.db.workspace_credentials.get", attributes)
	defer func() {
		span.End()
	}()

	wsc := &entities.WorkspaceCredentials{}

	transaction := database.SqlTxnFromContext(subctx, sql.client)
	tx := transaction.WithContext(subctx).Model(wsc).
		Where(fmt.Sprintf(`"%s".id = ?`, entities.TableWsc), id).
		First(wsc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return wsc, nil
}
