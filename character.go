package d1pg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) CreateCharacter(ctx context.Context, character d1.Character) (id int, err error) {
	query := "INSERT INTO d1.characters (account_id, gameserver_id, name, sex, class_id, color_1, color_2, color_3, alignment, alignment_enabled, xp, kamas, bonus_points, bonus_points_spell, honor, disgrace, stats, map_id, cell, direction, spells, mount_id, mounting)" +
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)" +
		" RETURNING id;"

	var color1 *d1typ.Color
	if character.Color1 != "" {
		color1 = &character.Color1
	}
	var color2 *d1typ.Color
	if character.Color1 != "" {
		color2 = &character.Color2
	}
	var color3 *d1typ.Color
	if character.Color1 != "" {
		color3 = &character.Color3
	}

	stats := fmt.Sprintf("%d,%d,%d,%d,%d,%d",
		character.Stats.Vitality,
		character.Stats.Wisdom,
		character.Stats.Strength,
		character.Stats.Intelligence,
		character.Stats.Chance,
		character.Stats.Agility,
	)

	spells := make([]string, len(character.Spells))
	for i, v := range character.Spells {
		spells[i] = fmt.Sprintf("%d~%d~%d", v.Id, v.Level, v.Position)
	}

	var mountId *int
	if character.MountId != 0 {
		mountId = &character.MountId
	}

	err = repoError(
		r.pool.QueryRow(ctx, query,
			character.AccountId, character.GameServerId, character.Name, character.Sex, character.ClassId,
			color1, color2, color3, character.Alignment, character.AlignmentEnabled, character.XP, character.Kamas,
			character.BonusPoints, character.BonusPointsSpell, character.Honor, character.Disgrace, stats,
			character.GameMapId, character.Cell, character.Direction, spells, mountId, character.Mounting,
		).Scan(&id),
	)
	return
}

func (r *Repo) UpdateCharacter(ctx context.Context, character d1.Character) error {
	query := "UPDATE d1.characters" +
		" SET account_id = $2, gameserver_id = $3, name = $4, sex = $5, class_id = $6, color_1 = $7, color_2 = $8, color_3 = $9, alignment = $10, alignment_enabled = $11, xp = $12, kamas = $13, bonus_points = $14, bonus_points_spell = $15, honor = $16, disgrace = $17, stats = $18, map_id = $19, cell = $20, direction = $21, spells = $22, mount_id = $23, mounting = $24" +
		" WHERE id = $1;"

	var color1 *d1typ.Color
	if character.Color1 != "" {
		color1 = &character.Color1
	}
	var color2 *d1typ.Color
	if character.Color1 != "" {
		color2 = &character.Color2
	}
	var color3 *d1typ.Color
	if character.Color1 != "" {
		color3 = &character.Color3
	}

	stats := fmt.Sprintf("%d,%d,%d,%d,%d,%d",
		character.Stats.Vitality,
		character.Stats.Wisdom,
		character.Stats.Strength,
		character.Stats.Intelligence,
		character.Stats.Chance,
		character.Stats.Agility,
	)

	spells := make([]string, len(character.Spells))
	for i, v := range character.Spells {
		spells[i] = fmt.Sprintf("%d~%d~%d", v.Id, v.Level, v.Position)
	}

	var mountId *int
	if character.MountId != 0 {
		mountId = &character.MountId
	}

	tag, err := r.pool.Exec(ctx, query, character.Id,
		character.AccountId, character.GameServerId, character.Name, character.Sex, character.ClassId,
		color1, color2, color3, character.Alignment, character.AlignmentEnabled, character.XP, character.Kamas,
		character.BonusPoints, character.BonusPointsSpell, character.Honor, character.Disgrace, stats,
		character.GameMapId, character.Cell, character.Direction, spells, mountId, character.Mounting)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return d1.ErrNotFound
	}

	return nil
}

func (r *Repo) DeleteCharacter(ctx context.Context, id int) error {
	query := "DELETE FROM d1.characters" +
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

func (r *Repo) AllCharacters(ctx context.Context) (map[int]d1.Character, error) {
	return r.characters(ctx, "")
}

func (r *Repo) AllCharactersByAccountId(ctx context.Context, accountId string) (map[int]d1.Character, error) {
	return r.characters(ctx, "account_id = $1;", accountId)
}

func (r *Repo) Characters(ctx context.Context, gameServerId int) (map[int]d1.Character, error) {
	return r.characters(ctx, "gameserver_id = $1", gameServerId)
}

func (r *Repo) CharactersByAccountId(ctx context.Context, gameServerId int, accountId string) (map[int]d1.Character, error) {
	return r.characters(ctx, "gameserver_id = $1 AND account_id = $2", gameServerId, accountId)
}

func (r *Repo) CharactersByGameMapId(ctx context.Context, gameServerId int, gameMapId int) (map[int]d1.Character, error) {
	return r.characters(ctx, "gameserver_id = $1 AND map_id = $2", gameServerId, gameMapId)
}

func (r *Repo) Character(ctx context.Context, id int) (d1.Character, error) {
	var char d1.Character

	chars, err := r.characters(ctx, "id = $1", id)
	if err != nil {
		return char, err
	}

	if len(chars) != 1 {
		return char, d1.ErrNotFound
	}

	for k := range chars {
		char = chars[k]
	}

	return char, nil
}

func (r *Repo) characters(ctx context.Context, conditions string, args ...interface{}) (map[int]d1.Character, error) {
	query := "SELECT id, account_id, gameserver_id, name, sex, class_id, color_1, color_2, color_3, alignment, alignment_enabled, xp, kamas, bonus_points, bonus_points_spell, honor, disgrace, stats, map_id, cell, direction, spells, mount_id, mounting" +
		" FROM d1.characters"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characters := make(map[int]d1.Character)
	for rows.Next() {
		var character d1.Character
		var color1 *d1typ.Color
		var color2 *d1typ.Color
		var color3 *d1typ.Color
		var statsStr string
		var spellsSli []string
		var mountId *int

		err = rows.Scan(&character.Id, &character.AccountId, &character.GameServerId, &character.Name,
			&character.Sex, &character.ClassId, &color1, &color2, &color3,
			&character.Alignment, &character.AlignmentEnabled, &character.XP, &character.Kamas, &character.BonusPoints,
			&character.BonusPointsSpell, &character.Honor, &character.Disgrace, &statsStr, &character.GameMapId,
			&character.Cell, &character.Direction, &spellsSli, &mountId, &character.Mounting)
		if err != nil {
			return nil, err
		}

		if color1 != nil {
			character.Color1 = *color1
		}
		if color2 != nil {
			character.Color2 = *color2
		}
		if color3 != nil {
			character.Color3 = *color3
		}

		stats, err := characterStats(statsStr)
		if err != nil {
			return nil, err
		}
		character.Stats = stats

		spells, err := characterSpells(spellsSli)
		if err != nil {
			return nil, err
		}
		character.Spells = spells

		if mountId != nil {
			character.MountId = *mountId
		}

		characters[character.Id] = character
	}

	return characters, nil
}

func characterStats(s string) (d1.CharacterStats, error) {
	var stats d1.CharacterStats

	sli := strings.Split(s, ",")
	if len(sli) < 6 {
		return stats, errors.New("malformed stats")
	}

	vitality, err := strconv.Atoi(sli[0])
	if err != nil {
		return stats, err
	}
	stats.Vitality = vitality

	wisdom, err := strconv.Atoi(sli[1])
	if err != nil {
		return stats, err
	}
	stats.Wisdom = wisdom

	strength, err := strconv.Atoi(sli[2])
	if err != nil {
		return stats, err
	}
	stats.Strength = strength

	intelligence, err := strconv.Atoi(sli[3])
	if err != nil {
		return stats, err
	}
	stats.Intelligence = intelligence

	chance, err := strconv.Atoi(sli[4])
	if err != nil {
		return stats, err
	}
	stats.Chance = chance

	agility, err := strconv.Atoi(sli[5])
	if err != nil {
		return stats, err
	}
	stats.Agility = agility

	return stats, nil
}

func characterSpells(sli []string) ([]d1.CharacterSpell, error) {
	spells := make([]d1.CharacterSpell, len(sli))
	for i, v := range sli {
		var spell d1.CharacterSpell

		sli := strings.Split(v, "~")
		if len(sli) != 3 {
			return nil, errors.New("invalid spells string")
		}

		id, err := strconv.Atoi(sli[0])
		if err != nil {
			return nil, err
		}
		spell.Id = id

		level, err := strconv.Atoi(sli[1])
		if err != nil {
			return nil, err
		}
		spell.Level = level

		position, err := strconv.Atoi(sli[2])
		if err != nil {
			return nil, err
		}
		spell.Position = position

		spells[i] = spell
	}
	return spells, nil
}
