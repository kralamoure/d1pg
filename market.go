package d1pg

import (
	"context"
	"strconv"
	"strings"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) Markets(ctx context.Context, gameServerId int) (markets map[string]d1.Market, err error) {
	query := "SELECT id, gameserver_id, quantity_1, quantity_2, quantity_3, types, fee, max_level, max_per_account, max_hours" +
		" FROM d1.markets" +
		" WHERE gameserver_id = $1;"

	rows, err := r.pool.Query(ctx, query, gameServerId)
	if err != nil {
		return
	}
	defer rows.Close()

	markets = make(map[string]d1.Market)
	for rows.Next() {
		var market d1.Market
		var itemTypes string

		err = rows.Scan(&market.Id, &market.GameServerId, &market.Quantity1, &market.Quantity2, &market.Quantity3,
			&itemTypes, &market.Fee, &market.MaxLevel, &market.MaxPerAccount, &market.MaxHours)
		if err != nil {
			return
		}

		if itemTypes != "" {
			sli := strings.Split(itemTypes, ",")
			market.Types = make([]d1typ.ItemType, len(sli))
			for i, v := range sli {
				itemType, err2 := strconv.Atoi(v)
				if err2 != nil {
					err = err2
					return
				}
				market.Types[i] = d1typ.ItemType(itemType)
			}
		}

		markets[market.Id] = market
	}
	return
}