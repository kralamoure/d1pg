package retropg

import (
	"context"

	"github.com/kralamoure/retro"
)

func (r *Db) NPCDialogs(ctx context.Context) (dialogs map[int]retro.NPCDialog, err error) {
	query := "SELECT id, text, responses" +
		" FROM d1_static.npc_dialogs;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	dialogs = make(map[int]retro.NPCDialog)
	for rows.Next() {
		var dialog retro.NPCDialog
		err = rows.Scan(&dialog.Id, &dialog.Text, &dialog.Responses)
		if err != nil {
			return
		}
		dialogs[dialog.Id] = dialog
	}
	return
}
