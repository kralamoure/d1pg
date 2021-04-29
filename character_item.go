package d1pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) CreateCharacterItem(ctx context.Context, item d1.CharacterItem) (id int, err error) {
	query := "INSERT INTO d1.characters_items (template_id, quantity, effects, position, character_id)" +
		" VALUES ($1, $2, $3, $4, $5)" +
		" RETURNING id;"

	effects := d1.EncodeItemEffects(item.Effects)

	var position *d1typ.CharacterItemPosition
	if item.Position != -1 {
		position = &item.Position
	}

	err = repoError(
		r.pool.QueryRow(ctx, query,
			item.TemplateId, item.Quantity, strings.Join(effects, ","), position, item.CharacterId,
		).Scan(&id),
	)
	return
}

func (r *Repo) UpdateCharacterItem(ctx context.Context, item d1.CharacterItem) error {
	query := "UPDATE d1.characters_items" +
		" SET template_id = $2, quantity = $3, effects = $4, position = $5, character_id = $6" +
		" WHERE id = $1;"

	effects := d1.EncodeItemEffects(item.Effects)

	var position *d1typ.CharacterItemPosition
	if item.Position != -1 {
		position = &item.Position
	}

	tag, err := r.pool.Exec(ctx, query, item.Id,
		item.TemplateId, item.Quantity, strings.Join(effects, ","), position, item.CharacterId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return d1.ErrNotFound
	}

	return nil
}

func (r *Repo) DeleteCharacterItem(ctx context.Context, id int) error {
	query := "DELETE FROM d1.characters_items" +
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

func (r *Repo) CharacterItemsByCharacterId(ctx context.Context, characterId int) (items map[int]d1.CharacterItem, err error) {
	return r.characterItems(ctx, "character_id = $1", characterId)
}

func (r *Repo) CharacterItem(ctx context.Context, id int) (d1.CharacterItem, error) {
	var item d1.CharacterItem

	items, err := r.characterItems(ctx, "id = $1", id)
	if err != nil {
		return item, err
	}

	if len(items) != 1 {
		return item, d1.ErrNotFound
	}

	for k := range items {
		item = items[k]
	}

	return item, nil
}

func (r *Repo) characterItems(ctx context.Context, conditions string, args ...interface{}) (map[int]d1.CharacterItem, error) {
	query := "SELECT id, template_id, quantity, effects, position, character_id" +
		" FROM d1.characters_items"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[int]d1.CharacterItem)
	for rows.Next() {
		var item d1.CharacterItem
		var effectsStr string
		var position *d1typ.CharacterItemPosition

		err = rows.Scan(&item.Id, &item.TemplateId, &item.Quantity, &effectsStr, &position, &item.CharacterId)
		if err != nil {
			return nil, err
		}

		if effectsStr != "" {
			effects, err := d1.DecodeItemEffects(strings.Split(effectsStr, ","))
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
