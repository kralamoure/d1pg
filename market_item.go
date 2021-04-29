package d1pg

import (
	"context"
	"strings"

	"github.com/kralamoure/d1"
)

func (r *Repo) CreateMarketItem(ctx context.Context, item d1.MarketItem) (id int, err error) {
	query := "INSERT INTO d1.markets_items (template_id, quantity, effects, price, market_id)" +
		" VALUES ($1, $2, $3, $4, $5)" +
		" RETURNING id;"

	effects := d1.EncodeItemEffects(item.Effects)

	err = repoError(
		r.pool.QueryRow(ctx, query,
			item.TemplateId, item.Quantity, strings.Join(effects, ","), item.Price, item.MarketId,
		).Scan(&id),
	)
	return
}

func (r *Repo) DeleteMarketItem(ctx context.Context, id int) error {
	query := "DELETE FROM d1.markets_items" +
		" WHERE id = $1;"

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return d1.ErrNotFound
	}
	return nil
}

func (r *Repo) MarketItems(ctx context.Context, gameServerId int) (items map[int]d1.MarketItem, err error) {
	query := "SELECT id, template_id, quantity, effects, price, market_id" +
		" FROM d1.markets_items" +
		" INNER JOIN markets m ON market_id = m.id" +
		" WHERE m.gameserver_id = $1;"

	rows, err := r.pool.Query(ctx, query, gameServerId)
	if err != nil {
		return
	}
	defer rows.Close()

	items = make(map[int]d1.MarketItem)
	for rows.Next() {
		var item d1.MarketItem
		var effectsStr string

		err = rows.Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &item.Price, &item.MarketId)
		if err != nil {
			return
		}

		if effectsStr != "" {
			effects, err2 := d1.DecodeItemEffects(strings.Split(effectsStr, ","))
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

func (r *Repo) MarketItemsByMarketId(ctx context.Context, marketId string) (items map[int]d1.MarketItem, err error) {
	query := "SELECT id, template_id, quantity, effects, price, market_id" +
		" FROM d1.markets_items" +
		" WHERE market_id = $1;"

	rows, err := r.pool.Query(ctx, query, marketId)
	if err != nil {
		return
	}
	defer rows.Close()

	items = make(map[int]d1.MarketItem)
	for rows.Next() {
		var item d1.MarketItem
		var effectsStr string

		err = rows.Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &item.Price, &item.MarketId)
		if err != nil {
			return
		}

		if effectsStr != "" {
			effects, err2 := d1.DecodeItemEffects(strings.Split(effectsStr, ","))
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

func (r *Repo) MarketItem(ctx context.Context, id int) (item d1.MarketItem, err error) {
	query := "SELECT id, template_id, quantity, effects, price, market_id" +
		" FROM d1.markets_items" +
		" WHERE market_id = $1;"

	var effectsStr string

	err = repoError(
		r.pool.QueryRow(ctx, query, id).Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &item.Price, &item.MarketId),
	)
	if err != nil {
		return
	}

	if effectsStr != "" {
		effects, err2 := d1.DecodeItemEffects(strings.Split(effectsStr, ","))
		if err2 != nil {
			err = err2
			return
		}
		item.Effects = effects
	}

	return
}
