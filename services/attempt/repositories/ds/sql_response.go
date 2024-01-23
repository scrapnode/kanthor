package ds

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Check(ctx context.Context, epId string, msgIds []string) (map[string][]int, error) {
	returning := map[string][]int{}
	if len(msgIds) == 0 {
		return returning, nil
	}

	conditions := []string{}
	values := map[string]any{}
	for i := range msgIds {
		// start with false value
		returning[msgIds[i]] = []int{}

		values[fmt.Sprintf("ep_id_%d", i)] = epId
		values[fmt.Sprintf("msg_id_%d", i)] = msgIds[i]
		conditions = append(conditions, fmt.Sprintf("(ep_id = @ep_id_%d AND msg_id = @msg_id_%d)", i, i))
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableRes)
	statement := fmt.Sprintf(
		"SELECT ep_id, msg_id, status FROM %s WHERE %s;",
		tableName,
		strings.Join(conditions, " OR "),
	)

	rows, err := sql.client.Raw(statement, values).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var epId, msgId string
		var status int
		if err := rows.Scan(&epId, &msgId, &status); err != nil {
			return nil, err
		}

		returning[msgId] = append(returning[msgId], status)
	}

	return returning, nil
}
