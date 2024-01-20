package ds

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Check(ctx context.Context, pairs []string) (map[string]bool, error) {
	returning := map[string]bool{}
	if len(pairs) == 0 {
		return returning, nil
	}

	conditions := []string{}
	values := map[string]any{}
	for i := range pairs {
		// start with false value
		returning[pairs[i]] = false

		ids := strings.Split(pairs[i], "/")
		values[fmt.Sprintf("@ep_id_%d", i)] = ids[0]
		values[fmt.Sprintf("@msg_id_%d", i)] = ids[1]
		conditions = append(conditions, fmt.Sprintf("(ep_id = @ep_id_%d AND msg_id = @msg_id_%d)", i, i))
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableReq)
	statement := fmt.Sprintf(
		"SELECT ep_id, msg_id FROM %s WHERE %s;",
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
		if err := rows.Scan(&epId, &msgId); err != nil {
			return nil, err
		}

		// then if we have any request, set the true value
		returning[fmt.Sprintf("%s/%s", epId, msgId)] = true
	}

	return returning, nil
}
