package retropg

import (
	"context"
	"strconv"
	"strings"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Db) NPCTemplates(ctx context.Context) (templates map[int]retro.NPCTemplate, err error) {
	query := "SELECT id, name, actions" +
		" FROM retro_static.npcs;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	templates = make(map[int]retro.NPCTemplate)
	for rows.Next() {
		var t retro.NPCTemplate
		var actions string
		err = rows.Scan(&t.Id, &t.Name, &actions)
		if err != nil {
			return
		}

		if actions != "" {
			sli := strings.Split(actions, ",")
			t.Actions = make([]retrotyp.NPCAction, len(sli))
			for i, v := range sli {
				action, err2 := strconv.Atoi(v)
				if err2 != nil {
					err = err2
					return
				}
				t.Actions[i] = retrotyp.NPCAction(action)
			}
		}

		templates[t.Id] = t
	}
	return
}
