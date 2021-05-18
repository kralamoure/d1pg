package retropg

import (
	"context"
	"fmt"
	"strings"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Storer) CreateCharacterItem(ctx context.Context, item retro.CharacterItem) (id int, err error) {
	query := "INSERT INTO retro.characters_items (template_id, quantity, effects, position, character_id)" +
		" VALUES ($1, $2, $3, $4, $5)" +
		" RETURNING id;"

	effects := retro.EncodeItemEffects(item.Effects)

	var position *retrotyp.CharacterItemPosition
	if item.Position != -1 {
		position = &item.Position
	}

	err = storerError(
		r.pool.QueryRow(ctx, query,
			item.TemplateId, item.Quantity, strings.Join(effects, ","), position, item.CharacterId,
		).Scan(&id),
	)
	return
}

func (r *Storer) UpdateCharacterItem(ctx context.Context, item retro.CharacterItem) error {
	query := "UPDATE retro.characters_items" +
		" SET template_id = $2, quantity = $3, effects = $4, position = $5, character_id = $6" +
		" WHERE id = $1;"

	effects := retro.EncodeItemEffects(item.Effects)

	var position *retrotyp.CharacterItemPosition
	if item.Position != -1 {
		position = &item.Position
	}

	tag, err := r.pool.Exec(ctx, query, item.Id,
		item.TemplateId, item.Quantity, strings.Join(effects, ","), position, item.CharacterId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return retro.ErrNotFound
	}

	return nil
}

func (r *Storer) DeleteCharacterItem(ctx context.Context, id int) error {
	query := "DELETE FROM retro.characters_items" +
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

func (r *Storer) CharacterItemsByCharacterId(ctx context.Context, characterId int) (items map[int]retro.CharacterItem, err error) {
	return r.characterItems(ctx, "character_id = $1", characterId)
}

func (r *Storer) CharacterItem(ctx context.Context, id int) (retro.CharacterItem, error) {
	var item retro.CharacterItem

	items, err := r.characterItems(ctx, "id = $1", id)
	if err != nil {
		return item, err
	}

	if len(items) != 1 {
		return item, retro.ErrNotFound
	}

	for k := range items {
		item = items[k]
	}

	return item, nil
}

func (r *Storer) characterItems(ctx context.Context, conditions string, args ...interface{}) (map[int]retro.CharacterItem, error) {
	query := "SELECT id, template_id, quantity, effects, position, character_id" +
		" FROM retro.characters_items"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[int]retro.CharacterItem)
	for rows.Next() {
		var item retro.CharacterItem
		var effectsStr string
		var position *retrotyp.CharacterItemPosition

		err = rows.Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &position, &item.CharacterId)
		if err != nil {
			return nil, err
		}

		if effectsStr != "" {
			effects, err := retro.DecodeItemEffects(strings.Split(effectsStr, ","))
			if err != nil {
				return nil, err
			}
			item.Effects = effects
		}

		if position == nil {
			item.Position = -1
		} else {
			item.Position = *position
		}

		items[item.Id] = item
	}

	return items, nil
}
