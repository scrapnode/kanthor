package ds

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Response, error) {
	doc := &entities.Response{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ?`, doc.TableName()), epId).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName())).
		Select([]string{
			fmt.Sprintf(`"%s"."id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."timestamp"`, doc.TableName()),
			fmt.Sprintf(`"%s"."ep_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."msg_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."req_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."app_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."type"`, doc.TableName()),
			fmt.Sprintf(`"%s"."uri"`, doc.TableName()),
			fmt.Sprintf(`"%s"."status"`, doc.TableName()),
			fmt.Sprintf(`"%s"."error"`, doc.TableName()),
		})

	tx = datastore.SqlApplyScanQuery(tx, entities.IdNsMsg, "msg_id", query)

	var docs []entities.Response
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlResponse) Get(ctx context.Context, epId, msgId string) ([]entities.Response, error) {

	doc := &entities.Response{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ? AND "%s"."msg_id" = ?`, doc.TableName(), doc.TableName()), epId, msgId).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName()))

	var docs []entities.Response
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}
