package ds

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/status"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) ListMessages(ctx context.Context, epId string, msgIds []string) (*MessageResponsetMaps, error) {
	doc := &entities.Response{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ? AND "%s"."msg_id" IN ?`, doc.TableName(), doc.TableName()), epId, msgIds).
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

	var responses []entities.Response
	if tx = tx.Find(&responses); tx.Error != nil {
		return nil, tx.Error
	}

	returning := &MessageResponsetMaps{Maps: make(map[string][]entities.Response), Success: make(map[string]string)}
	for _, response := range responses {
		if _, exist := returning.Maps[response.MsgId]; !exist {
			returning.Maps[response.MsgId] = []entities.Response{}
		}

		if status.Is2xx(response.Status) {
			returning.Success[response.MsgId] = response.MsgId
		}

		returning.Maps[response.MsgId] = append(returning.Maps[response.MsgId], response)
	}

	return returning, nil
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

	condition := &datastore.ScanningCondition{
		PrimaryKeyNs:  entities.IdNsMsg,
		PrimaryKeyCol: fmt.Sprintf(`"%s"."msg_id"`, doc.TableName()),
	}
	tx = datastore.SqlApplyScanQuery(tx, query, condition)

	var docs []entities.Response
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlResponse) Get(ctx context.Context, epId, msgId, id string) (*entities.Response, error) {
	doc := &entities.Response{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ? AND "%s"."msg_id" = ? AND "%s"."id" = ?`, doc.TableName(), doc.TableName(), doc.TableName()), epId, msgId, id).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName()))

	var docs []entities.Response
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
