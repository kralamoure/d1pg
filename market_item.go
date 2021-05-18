package retropg

import (
	"context"
	"strings"

	"github.com/kralamoure/retro"
)

func (r *Storer) CreateMarketItem(ctx context.Context, item retro.MarketItem) (id int, err error) {
	query := "INSERT INTO retro.markets_items (template_id, quantity, effects, price, market_id)" +
		" VALUES ($1, $2, $3, $4, $5)" +
		" RETURNING id;"

	effects := retro.EncodeItemEffects(item.Effects)

	err = storerError(
		r.pool.QueryRow(ctx, query,
			item.TemplateId, item.Quantity, strings.Join(effects, ","), item.Price, item.MarketId,
		).Scan(&id),
	)
	return
}

func (r *Storer) DeleteMarketItem(ctx context.Context, id int) error {
	query := "DELETE FROM retro.markets_items" +
		" WHERE id = $1;"

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return retro.ErrNotFound
	}
	return nil
}

func (r *Storer) MarketItems(ctx context.Context, gameServerId int) (items map[int]retro.MarketItem, err error) {
	query := "SELECT id, template_id, quantity, effects, price, market_id" +
		" FROM retro.markets_items" +
		" INNER JOIN markets m ON market_id = m.id" +
		" WHERE m.gameserver_id = $1;"

	rows, err := r.pool.Query(ctx, query, gameServerId)
	if err != nil {
		return
	}
	defer rows.Close()

	items = make(map[int]retro.MarketItem)
	for rows.Next() {
		var item retro.MarketItem
		var effectsStr string

		err = rows.Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &item.Price, &item.MarketId)
		if err != nil {
			return
		}

		if effectsStr != "" {
			effects, err2 := retro.DecodeItemEffects(strings.Split(effectsStr, ","))
			if err2 != nil {
				err = err2
				return
			}
			item.Effects = effects
		}

		items[item.Id] = item
	}

	return
}

func (r *Storer) MarketItemsByMarketId(ctx context.Context, marketId string) (items map[int]retro.MarketItem, err error) {
	query := "SELECT id, template_id, quantity, effects, price, market_id" +
		" FROM retro.markets_items" +
		" WHERE market_id = $1;"

	rows, err := r.pool.Query(ctx, query, marketId)
	if err != nil {
		return
	}
	defer rows.Close()

	items = make(map[int]retro.MarketItem)
	for rows.Next() {
		var item retro.MarketItem
		var effectsStr string

		err = rows.Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &item.Price, &item.MarketId)
		if err != nil {
			return
		}

		if effectsStr != "" {
			effects, err2 := retro.DecodeItemEffects(strings.Split(effectsStr, ","))
			if err2 != nil {
				err = err2
				return
			}
			item.Effects = effects
		}

		items[item.Id] = item
	}

	return
}

func (r *Storer) MarketItem(ctx context.Context, id int) (item retro.MarketItem, err error) {
	query := "SELECT id, template_id, quantity, effects, price, market_id" +
		" FROM retro.markets_items" +
		" WHERE market_id = $1;"

	var effectsStr string

	err = storerError(
		r.pool.QueryRow(ctx, query, id).Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &item.Price, &item.MarketId),
	)
	if err != nil {
		return
	}

	if effectsStr != "" {
		effects, err2 := retro.DecodeItemEffects(strings.Split(effectsStr, ","))
		if err2 != nil {
			err = err2
			return
		}
		item.Effects = effects
	}

	return
}
