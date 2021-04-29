package d1pg

import (
	"context"

	"github.com/kralamoure/d1"
)

func (r *Repo) EffectTemplates(ctx context.Context) (templates map[int]d1.EffectTemplate, err error) {
	query := "SELECT id, description, dice, operator, characteristic_id, element" +
		" FROM d1_static.effects;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	templates = make(map[int]d1.EffectTemplate)
	for rows.Next() {
		var t d1.EffectTemplate
		err = rows.Scan(&t.Id, &t.Description, &t.Dice, &t.Operator, &t.CharacteristicId, &t.Element)
		if err != nil {
			return
		}
		templates[t.Id] = t
	}
	return
}
