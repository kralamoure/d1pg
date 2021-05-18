package retropg

import (
	"context"
	"fmt"

	"github.com/kralamoure/retro"
)

func (r *Db) Triggers(ctx context.Context) (map[string]retro.Trigger, error) {
	return r.triggers(ctx, "")
}

func (r *Db) TriggerByGameMapIdAndCellId(ctx context.Context, gameMapId, cellId int) (retro.Trigger, error) {
	var trigger retro.Trigger

	triggers, err := r.triggers(ctx, "map_id = $1 AND cell_id = $2", gameMapId, cellId)
	if err != nil {
		return trigger, err
	}

	if len(triggers) != 1 {
		return trigger, retro.ErrNotFound
	}

	for k := range triggers {
		trigger = triggers[k]
	}

	return trigger, nil
}

func (r *Db) triggers(ctx context.Context, conditions string, args ...interface{}) (map[string]retro.Trigger, error) {
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

	triggers := make(map[string]retro.Trigger)
	for rows.Next() {
		var trigger retro.Trigger

		err = rows.Scan(&trigger.Id, &trigger.GameMapId, &trigger.CellId, &trigger.TargetGameMapId, &trigger.TargetCellId)
		if err != nil {
			return nil, err
		}

		triggers[trigger.Id] = trigger
	}

	return triggers, nil
}
