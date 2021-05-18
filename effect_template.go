package retropg

import (
	"context"

	"github.com/kralamoure/retro"
)

func (r *Db) EffectTemplates(ctx context.Context) (templates map[int]retro.EffectTemplate, err error) {
	query := "SELECT id, description, dice, operator, characteristic_id, element" +
		" FROM retro_static.effects;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	templates = make(map[int]retro.EffectTemplate)
	for rows.Next() {
		var t retro.EffectTemplate
		err = rows.Scan(&t.Id, &t.Description, &t.Dice, &t.Operator, &t.CharacteristicId, &t.Element)
		if err != nil {
			return
		}
		templates[t.Id] = t
	}
	return
}
