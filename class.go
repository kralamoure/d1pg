package retropg

import (
	"context"
	"encoding/json"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Storer) Classes(ctx context.Context) (classes map[retrotyp.ClassId]retro.Class, err error) {
	query := "SELECT id, name, label, short_description, description, spells, boost_costs" +
		" FROM d1_static.classes;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	classes = make(map[retrotyp.ClassId]retro.Class)
	for rows.Next() {
		var class retro.Class
		var boostCostsStr string

		err = rows.Scan(&class.Id, &class.Name, &class.Label, &class.ShortDescription, &class.Description,
			&class.Spells, &boostCostsStr)
		if err != nil {
			return
		}

		var boostCosts [][][]int

		err := json.Unmarshal([]byte(boostCostsStr), &boostCosts)
		if err != nil {
			return nil, err
		}

		for i, v := range boostCosts {
			characteristic := make([]retro.ClassBoostCost, len(v))

			for i, v := range v {
				var cost retro.ClassBoostCost
				cost.Quantity = v[0]
				cost.Cost = v[1]
				if len(v) >= 3 {
					cost.Bonus = v[2]
				} else {
					cost.Bonus = 1
				}
				characteristic[i] = cost
			}

			switch i {
			case 0:
				class.BoostCosts.Vitality = characteristic
			case 1:
				class.BoostCosts.Wisdom = characteristic
			case 2:
				class.BoostCosts.Strength = characteristic
			case 3:
				class.BoostCosts.Intelligence = characteristic
			case 4:
				class.BoostCosts.Chance = characteristic
			case 5:
				class.BoostCosts.Agility = characteristic
			}
		}

		classes[class.Id] = class
	}
	return
}
