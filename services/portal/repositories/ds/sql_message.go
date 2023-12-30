package ds

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) ListByIds(ctx context.Context, appId string, ids []string) ([]entities.Message, error) {
	doc := &entities.Message{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."app_id" = ? AND "%s"."id" IN ?`, doc.TableName(), doc.TableName()), appId, ids).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName())).
		Select([]string{
			fmt.Sprintf(`"%s"."id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."timestamp"`, doc.TableName()),
			fmt.Sprintf(`"%s"."app_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."type"`, doc.TableName()),
		})

	var docs []entities.Message
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlMessage) GetByIds(ctx context.Context, appId string, ids []string) ([]entities.Message, error) {
	doc := &entities.Message{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."app_id" = ? AND "%s"."id" IN ?`, doc.TableName(), doc.TableName()), appId, ids).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName()))

	var docs []entities.Message
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}
func (sql *SqlMessage) List(ctx context.Context, appId string, query *entities.ScanningQuery) ([]entities.Message, error) {
	doc := &entities.Message{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."app_id" = ?`, doc.TableName()), appId).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName())).
		Select([]string{
			fmt.Sprintf(`"%s"."id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."timestamp"`, doc.TableName()),
			fmt.Sprintf(`"%s"."app_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."type"`, doc.TableName()),
		})
	condition := &datastore.ScanningCondition{
		PrimaryKeyNs:  entities.IdNsMsg,
		PrimaryKeyCol: fmt.Sprintf(`"%s"."id"`, doc.TableName()),
	}
	tx = datastore.SqlApplyScanQuery(tx, query, condition)

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
