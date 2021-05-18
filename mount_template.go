package retropg

import (
	"context"
	"fmt"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Storer) MountTemplates(ctx context.Context) (map[int]retro.MountTemplate, error) {
	return r.mountTemplates(ctx, "")
}

func (r *Storer) mountTemplates(ctx context.Context, conditions string, args ...interface{}) (map[int]retro.MountTemplate, error) {
	query := "SELECT id, name, gfx_id, color_1, color_2, color_3, max_effects" +
		" FROM d1_static.mounts"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mountTemplates := make(map[int]retro.MountTemplate)
	for rows.Next() {
		var mountTemplate retro.MountTemplate
		var color1 *retrotyp.Color
		var color2 *retrotyp.Color
		var color3 *retrotyp.Color
		var maxEffectsSli []string

		err = rows.Scan(&mountTemplate.Id, &mountTemplate.Name, &mountTemplate.GFXId, &color1, &color2, &color3, &maxEffectsSli)
		if err != nil {
			return nil, err
		}

		if color1 != nil {
			mountTemplate.Color1 = *color1
		}
		if color2 != nil {
			mountTemplate.Color2 = *color2
		}
		if color3 != nil {
			mountTemplate.Color3 = *color3
		}

		maxEffects, err := retro.DecodeItemEffects(maxEffectsSli)
		if err != nil {
			return nil, err
		}
		mountTemplate.MaxEffects = maxEffects

		mountTemplates[mountTemplate.Id] = mountTemplate
	}

	return mountTemplates, nil
}
