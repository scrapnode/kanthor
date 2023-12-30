package ds

import (
	"context"
	"fmt"

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
		sql.mapping(returning, &response)
	}
	return returning, nil
}

func (sql *SqlResponse) GetMessages(ctx context.Context, epId string, msgIds []string) (*MessageResponsetMaps, error) {
	doc := &entities.Response{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."ep_id" = ? AND "%s"."msg_id" IN ?`, doc.TableName(), doc.TableName()), epId, msgIds).
		Order(fmt.Sprintf(`"%s"."app_id" DESC, "%s"."msg_id" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName()))

	var responses []entities.Response
	if tx = tx.Find(&responses); tx.Error != nil {
		return nil, tx.Error
	}

	returning := &MessageResponsetMaps{Maps: make(map[string][]entities.Response), Success: make(map[string]string)}
	for _, response := range responses {
		sql.mapping(returning, &response)
	}
	return returning, nil
}

func (sql *SqlResponse) mapping(returning *MessageResponsetMaps, response *entities.Response) {
	if _, exist := returning.Maps[response.MsgId]; !exist {
		returning.Maps[response.MsgId] = []entities.Response{}
	}

	if status.Is2xx(response.Status) {
		returning.Success[response.MsgId] = response.MsgId
	}

	returning.Maps[response.MsgId] = append(returning.Maps[response.MsgId], *response)
}
