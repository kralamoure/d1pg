package d1pg

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
	"github.com/kralamoure/d1encoding"
)

func (r *Repo) Spells(ctx context.Context) (spells map[int]d1.Spell, err error) {
	query := "SELECT id, name, description, levels" +
		" FROM d1_static.spells;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	spells = make(map[int]d1.Spell)
	for rows.Next() {
		var spell d1.Spell
		var levels []string

		err = rows.Scan(&spell.Id, &spell.Name, &spell.Description, &levels)
		if err != nil {
			return
		}

		for i, v := range levels {
			level, err2 := decodeSpellLevel(v, i+1)
			if err2 != nil {
				err = err2
				return
			}
			spell.Levels = append(spell.Levels, level)
		}

		spells[spell.Id] = spell
	}
	return
}

func decodeSpellLevel(s string, grade int) (level d1typ.SpellLevel, err error) {
	var sli []interface{}
	err = json.Unmarshal([]byte(s), &sli)
	if err != nil {
		return
	}
	if len(sli) != 20 {
		err = errors.New("invalid length for spell level slice")
		return
	}

	level.Grade = grade

	if sli[0] != nil {
		effects, err2 := decodeSpellLevelEffects(sli[0])
		if err2 != nil {
			err = err2
			return
		}
		level.Effects = effects
	}

	if sli[1] != nil {
		effectsCritical, err2 := decodeSpellLevelEffects(sli[1])
		if err2 != nil {
			err = err2
			return
		}
		level.EffectsCritical = effectsCritical
	}

	apCost, ok := sli[2].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.APCost = int(apCost)

	rangeN, ok := sli[3].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.Range = int(rangeN)

	rangeMax, ok := sli[4].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.RangeMax = int(rangeMax)

	criticalHitProbability, ok := sli[5].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.CriticalHitProbability = int(criticalHitProbability)

	criticalFailureProbability, ok := sli[6].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.CriticalFailureProbability = int(criticalFailureProbability)

	level.Linear, ok = sli[7].(bool)
	if !ok {
		err = errInvalidAssertion
		return
	}

	level.RequiresLineOfSight, ok = sli[8].(bool)
	if !ok {
		err = errInvalidAssertion
		return
	}

	level.RequiresFreeCell, ok = sli[9].(bool)
	if !ok {
		err = errInvalidAssertion
		return
	}

	level.AdjustableRange, ok = sli[10].(bool)
	if !ok {
		err = errInvalidAssertion
		return
	}

	classId, ok := sli[11].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.ClassId = d1typ.ClassId(classId)

	maxCastsPerTurn, ok := sli[12].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.MaxCastsPerTurn = int(maxCastsPerTurn)

	maxCastsPerTarget, ok := sli[13].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.MaxCastsPerTarget = int(maxCastsPerTarget)

	minCastInterval, ok := sli[14].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.MinCastInterval = int(minCastInterval)

	effectZones, ok := sli[15].(string)
	if !ok {
		err = errInvalidAssertion
		return
	}
	if effectZones != "" && len(effectZones)%2 == 0 {
		shapes := make([]d1typ.EffectZoneShape, len(effectZones)/2)
		sizes := make([]int, len(effectZones)/2)
		for i := 0; i < len(effectZones)/2; i++ {
			shape := d1typ.EffectZoneShape(effectZones[i*2])

			sizeR := rune(effectZones[i*2+1])
			size, err2 := d1encoding.Decode64(sizeR)
			if err2 != nil {
				err = err2
				return
			}

			shapes[i] = shape
			sizes[i] = size
		}

		if len(shapes) != len(level.Effects)+len(level.EffectsCritical) {
			err = errors.New("length of effect zones doesn't coincide with length of effects and critical effects")
			return
		}

		for i := range level.Effects {
			level.Effects[i].ZoneShape = shapes[i]
			level.Effects[i].ZoneSize = sizes[i]
		}

		for i := range level.EffectsCritical {
			offset := len(level.Effects)
			level.EffectsCritical[i].ZoneShape = shapes[i+offset]
			level.EffectsCritical[i].ZoneSize = sizes[i+offset]
		}
	}

	statesRequiredSli, ok := sli[16].([]interface{})
	if !ok {
		err = errInvalidAssertion
		return
	}
	statesRequired := make([]int, len(statesRequiredSli))
	for i, v := range statesRequiredSli {
		n, ok := v.(float64)
		if !ok {
			err = errInvalidAssertion
			return
		}
		statesRequired[i] = int(n)
	}
	level.StatesRequired = statesRequired

	statesForbiddenSli, ok := sli[17].([]interface{})
	if !ok {
		err = errInvalidAssertion
		return
	}
	statesForbidden := make([]int, len(statesForbiddenSli))
	for i, v := range statesForbiddenSli {
		n, ok := v.(float64)
		if !ok {
			err = errInvalidAssertion
			return
		}
		statesForbidden[i] = int(n)
	}
	level.StatesForbidden = statesForbidden

	minPlayerLevel, ok := sli[18].(float64)
	if !ok {
		err = errInvalidAssertion
		return
	}
	level.MinPlayerLevel = int(minPlayerLevel)

	level.CriticalFailureEndsTurn, ok = sli[19].(bool)
	if !ok {
		err = errInvalidAssertion
		return
	}

	return
}

func decodeSpellLevelEffects(v interface{}) (effects []d1typ.Effect, err error) {
	sli, ok := v.([]interface{})
	if !ok {
		err = errInvalidAssertion
		return
	}

	effects = make([]d1typ.Effect, len(sli))
	for i, v := range sli {
		effectSli, ok := v.([]interface{})
		if !ok {
			err = errInvalidAssertion
			return
		}

		if len(effectSli) < 7 {
			err = errors.New("invalid length for effect slice")
			return
		}

		var effect d1typ.Effect

		id, ok := effectSli[0].(float64)
		if !ok {
			err = errInvalidAssertion
			return
		}
		effect.Id = int(id)

		if effectSli[1] == nil {
			effect.DiceNum = -1
		} else {
			diceNum, _ := effectSli[1].(float64)
			effect.DiceNum = int(diceNum)
		}

		if effectSli[2] == nil {
			effect.DiceSide = -1
		} else {
			diceSide, _ := effectSli[2].(float64)
			effect.DiceSide = int(diceSide)
		}

		if effectSli[3] == nil {
			effect.Value = -1
		} else {
			value, _ := effectSli[3].(float64)
			effect.Value = int(value)
		}

		duration, ok := effectSli[4].(float64)
		if !ok {
			err = errInvalidAssertion
			return
		}
		effect.Duration = int(duration)

		random, ok := effectSli[5].(float64)
		if !ok {
			err = errInvalidAssertion
			return
		}
		effect.Random = int(random)

		targetId, ok := effectSli[6].(float64)
		if !ok {
			err = errInvalidAssertion
			return
		}
		effect.TargetId = int(targetId)

		if len(effectSli) >= 8 {
			effect.Param, _ = effectSli[7].(string)
		}

		effects[i] = effect
	}

	return
}
