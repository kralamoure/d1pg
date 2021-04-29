package d1pg

import (
	"context"
	"strconv"
	"strings"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) NPCTemplates(ctx context.Context) (templates map[int]d1.NPCTemplate, err error) {
	query := "SELECT id, name, actions" +
		" FROM d1_static.npcs;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	templates = make(map[int]d1.NPCTemplate)
	for rows.Next() {
		var t d1.NPCTemplate
		var actions string
		err = rows.Scan(&t.Id, &t.Name, &actions)
		if err != nil {
			return
		}

		if actions != "" {
			sli := strings.Split(actions, ",")
			t.Actions = make([]d1typ.NPCAction, len(sli))
			for i, v := range sli {
				action, err2 := strconv.Atoi(v)
				if err2 != nil {
					err = err2
					return
				}
				t.Actions[i] = d1typ.NPCAction(action)
			}
		}

		templates[t.Id] = t
	}
	return
}
