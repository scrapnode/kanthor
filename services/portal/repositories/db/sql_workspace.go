package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
}

func (sql *SqlWorkspace) Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	if err := doc.Validate(); err != nil {
		return nil, err
	}

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspace) Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	if err := doc.Validate(); err != nil {
		return nil, err
	}

	transaction := database.SqlTxnFromContext(ctx, sql.client)

	updates := []string{
		"name",
		"updated_at",
	}
	tx := transaction.
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		// When update with struct, GORM will only update non-zero fields,
		// you might want to use map to update attributes or use Select to specify fields to update
		Select(updates).
		Updates(doc)
	return doc, tx.Error
}

func (sql *SqlWorkspace) ListByIds(ctx context.Context, ids []string) ([]entities.Workspace, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var docs []entities.Workspace

	tx := sql.client.WithContext(ctx).
		Model(&entities.Workspace{}).
		Where("id IN ?", ids).
		Find(&docs)

	return docs, tx.Error
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}
	transaction := database.SqlTxnFromContext(ctx, sql.client)

	tx := transaction.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspace) ListOwned(ctx context.Context, owner string) ([]entities.Workspace, error) {
	var docs []entities.Workspace

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&entities.Workspace{}).
		Where("owner_id = ?", owner).
		Order("id DESC").
		Find(&docs)

	return docs, tx.Error
}

type SnaptshotRow struct {
	WsId                   string `json:"ws_id"`
	WsName                 string `json:"ws_name"`
	AppId                  string `json:"app_id"`
	AppName                string `json:"app_name"`
	EpId                   string `json:"ep_id"`
	EpName                 string `json:"ep_name"`
	EpMethod               string `json:"ep_method"`
	EpUri                  string `json:"ep_uri"`
	EprId                  string `json:"epr_id"`
	EprName                string `json:"epr_name"`
	EprPriority            int32  `json:"epr_priority"`
	EprExclusionary        bool   `json:"epr_exclusionary"`
	EprConditionSource     string `json:"epr_condition_source"`
	EprConditionExpression string `json:"epr_condition_expression"`
}

func (sql *SqlWorkspace) GetSnapshotRows(ct context.Context, id string) ([]SnaptshotRow, error) {
	// JOIN WITH CAUTION
	// Most of the time, you may prefer not to join more than two tables simultaneously due to the size of your dataset.
	// However, in this particular case, even in the worst scenario,
	// we may only have thousands of records, so there's no need to worry about joining multiple tables.
	raw := `
	SELECT
		kw.id AS ws_id,
		kw.name AS ws_name,
		ka.id AS app_id,
		ka.name AS app_name,
		ke.id AS ep_id,
		ke.name AS ep_name,
		ke.method AS ep_method,
		ke.uri AS ep_uri,
		ker.id AS epr_id,
		ker.name AS epr_name,
		ker.priority AS epr_priority,
		ker.exclusionary AS epr_exclusionary,
		ker.condition_source AS epr_condition_source,
		ker.condition_expression AS epr_condition_expression
	FROM %s kw 
	JOIN %s ka ON ka.ws_id = kw.id
	JOIN %s ke ON ke.app_id = ka.id
	JOIN %s ker ON ker.ep_id = ke.id
	WHERE kw.id = ?
	ORDER BY ka.id DESC, ke.id DESC, ker.id DESC 
	`
	query := fmt.Sprintf(raw, entities.TableWs, entities.TableApp, entities.TableEp, entities.TableEpr)

	rows := []SnaptshotRow{}
	tx := sql.client.Raw(query, id).Scan(&rows)
	return rows, tx.Error
}
