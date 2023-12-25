package ds

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlMessage) List(ctx context.Context, appId string, query *entities.ScanningQuery) ([]entities.Message, error) {
	doc := &entities.Message{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."app_id" = ?`, doc.TableName()), appId).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName()))

	tx = datastore.SqlApplyScanQuery(tx, sql.timer, entities.IdNsMsg, "id", query)

	var docs []entities.Message
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlMessage) Get(ctx context.Context, appId, id string) (*entities.Message, error) {
	doc := &entities.Message{}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."app_id" = ? AND "%s"."id" = ?`, doc.TableName(), doc.TableName()), appId, doc.Id).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName())).
		First(doc)
	if tx.Error != nil {
		return nil, datastore.SqlError(tx.Error)
	}

	return doc, nil
}
