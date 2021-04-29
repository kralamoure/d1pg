package d1pg

import (
	"context"
	"fmt"

	"github.com/kralamoure/d1"
)

func (r *Repo) Triggers(ctx context.Context) (map[string]d1.Trigger, error) {
	return r.triggers(ctx, "")
}

func (r *Repo) TriggerByGameMapIdAndCellId(ctx context.Context, gameMapId, cellId int) (d1.Trigger, error) {
	var trigger d1.Trigger

	triggers, err := r.triggers(ctx, "map_id = $1 AND cell_id = $2", gameMapId, cellId)
	if err != nil {
		return trigger, err
	}

	if len(triggers) != 1 {
		return trigger, d1.ErrNotFound
	}

	for k := range triggers {
		trigger = triggers[k]
	}

	return trigger, nil
}

func (r *Repo) triggers(ctx context.Context, conditions string, args ...interface{}) (map[string]d1.Trigger, error) {
	query := "SELECT id, map_id, cell_id, target_map_id, target_cell_id" +
		" FROM d1_static.triggers"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	triggers := make(map[string]d1.Trigger)
	for rows.Next() {
		var trigger d1.Trigger

		err = rows.Scan(&trigger.Id, &trigger.GameMapId, &trigger.CellId, &trigger.TargetGameMapId, &trigger.TargetCellId)
		if err != nil {
			return nil, err
		}

		triggers[trigger.Id] = trigger
	}

	return triggers, nil
}
