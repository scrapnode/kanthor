package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceExportIn struct {
	Id string
}

func (in *WorkspaceExportIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceExportOut struct {
	Data *entities.WorkspaceSnapshot
}

func (uc *workspace) Export(ctx context.Context, in *WorkspaceExportIn) (*WorkspaceExportOut, error) {
	rows, err := uc.repositories.Database().Workspace().GetSnapshotRows(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	snapshot := &entities.WorkspaceSnapshot{Applications: make(map[string]entities.WorkspaceSnapshotApp)}
	for _, row := range rows {
		snapshot.Id = row.WsId
		snapshot.Name = row.WsName

		// app mapping
		if _, exist := snapshot.Applications[row.AppId]; !exist {
			snapshot.Applications[row.AppId] = entities.WorkspaceSnapshotApp{
				Name:      row.AppName,
				Endpoints: make(map[string]entities.WorkspaceSnapshotEp),
			}
		}

		// ep mapping
		if _, exist := snapshot.Applications[row.AppId].Endpoints[row.EpId]; !exist {
			snapshot.Applications[row.AppId].Endpoints[row.EpId] = entities.WorkspaceSnapshotEp{
				Name:   row.EpName,
				Method: row.EpMethod,
				Uri:    row.EpUri,
				Rules:  make(map[string]entities.WorkspaceSnapshotEpr),
			}
		}

		// epr mapping
		if _, exist := snapshot.Applications[row.AppId].Endpoints[row.EpId].Rules[row.EprId]; !exist {
			snapshot.Applications[row.AppId].Endpoints[row.EpId].Rules[row.EprId] = entities.WorkspaceSnapshotEpr{
				Name:                row.EprName,
				Priority:            row.EprPriority,
				Exclusionary:        row.EprExclusionary,
				ConditionSource:     row.EprConditionSource,
				ConditionExpression: row.EprConditionExpression,
			}
		}
	}

	// IMPORTANT
	// Because the dataset is too small so we don't need to upload it to S3
	out := &WorkspaceExportOut{Data: snapshot}
	return out, nil
}
