package retropg

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Db) ItemSets(ctx context.Context) (itemSets map[int]retro.ItemSet, err error) {
	query := "SELECT id, name, bonus" +
		" FROM retro_static.itemsets;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	itemSets = make(map[int]retro.ItemSet)
	for rows.Next() {
		var itemSet retro.ItemSet
		var bonus string
		err = rows.Scan(&itemSet.Id, &itemSet.Name, &bonus)
		if err != nil {
			return
		}

		sli := strings.Split(bonus, ";")
		itemSet.Bonus = make([][]retrotyp.Effect, len(sli))
		for i, v := range sli {
			if v == "" {
				continue
			}
			stats := strings.Split(v, ",")
			itemSet.Bonus[i] = make([]retrotyp.Effect, len(stats))
			for i2, v := range stats {
				effectStr := strings.Split(v, ":")
				if len(effectStr) != 2 {
					err = errors.New("malformed bonus")
					return
				}
				id, err2 := strconv.Atoi(effectStr[0])
				if err2 != nil {
					err = err2
					return
				}
				diceNum, err2 := strconv.Atoi(effectStr[1])
				if err2 != nil {
					err = err2
					return
				}
				itemSet.Bonus[i][i2] = retrotyp.Effect{
					Id:        id,
					ZoneShape: retrotyp.EffectZoneShapeCircle,
					DiceNum:   diceNum,
					Hidden:    true,
				}
			}
		}
		itemSets[itemSet.Id] = itemSet
	}
	return
}
