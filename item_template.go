package retropg

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kralamoure/retro"
)

func (r *Db) ItemTemplates(ctx context.Context) (templates map[int]retro.ItemTemplate, err error) {
	query := "SELECT id, name, description, type, enhanceable, two_hands, ethereal, hidden, itemset_id, can_use, can_target, level, gfx, price, weight, cursed, conditions, weapon_effects, effects" +
		" FROM retro_static.items;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	templates = make(map[int]retro.ItemTemplate)
	for rows.Next() {
		var t retro.ItemTemplate
		var itemSetId *int
		var weaponEffectsStr string
		var effectsStr string
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Type,
			&t.Enhanceable, &t.TwoHands, &t.Ethereal, &t.Hidden,
			&itemSetId, &t.CanUse, &t.CanTarget, &t.Level,
			&t.GFX, &t.Price, &t.Weight, &t.Cursed,
			&t.Conditions, &weaponEffectsStr, &effectsStr)
		if err != nil {
			return
		}
		if itemSetId != nil {
			t.ItemSetId = *itemSetId
		}

		if weaponEffectsStr != "" {
			sli := strings.Split(weaponEffectsStr, ",")
			if len(sli) != 8 {
				err = errors.New("malformed weapon effects")
				return
			}
			t.WeaponEffects.CriticalHitBonus, err = strconv.Atoi(sli[0])
			if err != nil {
				return
			}
			t.WeaponEffects.APCost, err = strconv.Atoi(sli[1])
			if err != nil {
				return
			}
			t.WeaponEffects.RangeMin, err = strconv.Atoi(sli[2])
			if err != nil {
				return
			}
			t.WeaponEffects.RangeMax, err = strconv.Atoi(sli[3])
			if err != nil {
				return
			}
			t.WeaponEffects.CriticalHit, err = strconv.Atoi(sli[4])
			if err != nil {
				return
			}
			t.WeaponEffects.CriticalFailure, err = strconv.Atoi(sli[5])
			if err != nil {
				return
			}
			t.WeaponEffects.LineOnly, err = strconv.ParseBool(sli[6])
			if err != nil {
				return
			}
			t.WeaponEffects.LineOfSight, err = strconv.ParseBool(sli[7])
			if err != nil {
				return
			}
		}

		if effectsStr != "" {
			effects, err2 := retro.DecodeItemEffects(strings.Split(effectsStr, ","))
			if err2 != nil {
				err = err2
				return
			}
			t.Effects = effects
		}

		templates[t.Id] = t
	}
	return
}
