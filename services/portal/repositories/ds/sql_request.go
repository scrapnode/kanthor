package ds

import (
	"context"
	"fmt"
	"log"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Request, error) {
	log.Println("-----------------------------")
	doc := &entities.Request{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ?`, doc.TableName()), epId).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName())).
		Select([]string{
			fmt.Sprintf(`"%s"."id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."timestamp"`, doc.TableName()),
			fmt.Sprintf(`"%s"."ep_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."msg_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."app_id"`, doc.TableName()),
			fmt.Sprintf(`"%s"."type"`, doc.TableName()),
			fmt.Sprintf(`"%s"."uri"`, doc.TableName()),
			fmt.Sprintf(`"%s"."method"`, doc.TableName()),
		})

	tx = datastore.SqlApplyScanQuery(tx, entities.IdNsMsg, "msg_id", query)

	var docs []entities.Request
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlRequest) Get(ctx context.Context, epId, msgId, id string) (*entities.Request, error) {
	doc := &entities.Request{
		EpId:  epId,
		MsgId: msgId,
	}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ? AND "%s"."msg_id" = ? AND "%s"."id" = ?`, doc.TableName(), doc.TableName(), doc.TableName()), epId, msgId, id).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName()))

	if tx = tx.First(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
