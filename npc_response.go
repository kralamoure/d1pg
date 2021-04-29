package d1pg

import (
	"context"

	"github.com/kralamoure/d1"
)

func (r *Repo) NPCResponses(ctx context.Context) (responses map[int]d1.NPCResponse, err error) {
	query := "SELECT id, text, action, arguments, conditions" +
		" FROM d1_static.npc_responses;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	responses = make(map[int]d1.NPCResponse)
	for rows.Next() {
		var response d1.NPCResponse
		err = rows.Scan(&response.Id, &response.Text, &response.Action, &response.Arguments, &response.Conditions)
		if err != nil {
			return
		}
		responses[response.Id] = response
	}
	return
}
