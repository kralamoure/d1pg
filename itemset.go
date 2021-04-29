package d1pg

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) ItemSets(ctx context.Context) (itemSets map[int]d1.ItemSet, err error) {
	query := "SELECT id, name, bonus" +
		" FROM d1_static.itemsets;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	itemSets = make(map[int]d1.ItemSet)
	for rows.Next() {
		var itemSet d1.ItemSet
		var bonus string
		err = rows.Scan(&itemSet.Id, &itemSet.Name, &bonus)
		if err != nil {
			return
		}

		sli := strings.Split(bonus, ";")
		itemSet.Bonus = make([][]d1typ.Effect, len(sli))
		for i, v := range sli {
			if v == "" {
				continue
			}
			stats := strings.Split(v, ",")
			itemSet.Bonus[i] = make([]d1typ.Effect, len(stats))
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
				itemSet.Bonus[i][i2] = d1typ.Effect{
					Id:        id,
					ZoneShape: d1typ.EffectZoneShapeCircle,
					DiceNum:   diceNum,
					Hidden:    true,
				}
			}
		}
		itemSets[itemSet.Id] = itemSet
	}
	return
}
