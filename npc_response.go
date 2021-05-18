package retropg

import (
	"context"

	"github.com/kralamoure/retro"
)

func (r *Db) NPCResponses(ctx context.Context) (responses map[int]retro.NPCResponse, err error) {
	query := "SELECT id, text, action, arguments, conditions" +
		" FROM retro_static.npc_responses;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	responses = make(map[int]retro.NPCResponse)
	for rows.Next() {
		var response retro.NPCResponse
		err = rows.Scan(&response.Id, &response.Text, &response.Action, &response.Arguments, &response.Conditions)
		if err != nil {
			return
		}
		responses[response.Id] = response
	}
	return
}
