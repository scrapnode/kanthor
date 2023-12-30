package ds

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) ListMessages(ctx context.Context, epId string, query *entities.ScanningQuery) (*MessageRequestMaps, error) {
	returning := &MessageRequestMaps{Maps: make(map[string][]entities.Request)}

	for {
		requests, err := sql.List(ctx, epId, query)
		if err != nil {
			return nil, err
		}

		// no more data
		if len(requests) == 0 {
			return returning, nil
		}

		for _, request := range requests {
			if _, exist := returning.Maps[request.MsgId]; !exist {
				returning.Maps[request.MsgId] = []entities.Request{}
				returning.MsgIds = append(returning.MsgIds, request.MsgId)
			}
			returning.Maps[request.MsgId] = append(returning.Maps[request.MsgId], request)

			// checking limit exceeded
			if len(returning.Maps) >= query.Limit {
				return returning, nil
			}
		}

		// checking no more data
		if len(requests) < query.Limit {
			return returning, nil
		}

		// checking limit exceeded
		if len(returning.Maps) >= query.Limit {
			return returning, nil
		}
	}
}

func (sql *SqlRequest) GetMessage(ctx context.Context, epId, msgId string) (*MessageRequestMaps, error) {
	doc := &entities.Request{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ? AND "%s"."msg_id" = ?`, doc.TableName(), doc.TableName()), epId, msgId).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName()))

	var requests []entities.Request
	if tx = tx.Find(&requests); tx.Error != nil {
		return nil, tx.Error
	}

	returning := &MessageRequestMaps{Maps: make(map[string][]entities.Request)}
	for _, request := range requests {
		if _, exist := returning.Maps[request.MsgId]; !exist {
			returning.Maps[request.MsgId] = []entities.Request{}
			returning.MsgIds = append(returning.MsgIds, request.MsgId)
		}
		returning.Maps[request.MsgId] = append(returning.Maps[request.MsgId], request)
	}

	return returning, nil
}

func (sql *SqlRequest) List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Request, error) {
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

	condition := &datastore.ScanningCondition{
		PrimaryKeyNs:  entities.IdNsMsg,
		PrimaryKeyCol: fmt.Sprintf(`"%s"."msg_id"`, doc.TableName()),
	}
	tx = datastore.SqlApplyScanQuery(tx, query, condition)

	var docs []entities.Request
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}
