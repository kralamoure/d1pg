package d1pg

import (
	"context"

	"github.com/kralamoure/d1"
)

func (r *Repo) NPCDialogs(ctx context.Context) (dialogs map[int]d1.NPCDialog, err error) {
	query := "SELECT id, text, responses" +
		" FROM d1_static.npc_dialogs;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	dialogs = make(map[int]d1.NPCDialog)
	for rows.Next() {
		var dialog d1.NPCDialog
		err = rows.Scan(&dialog.Id, &dialog.Text, &dialog.Responses)
		if err != nil {
			return
		}
		dialogs[dialog.Id] = dialog
	}
	return
}
