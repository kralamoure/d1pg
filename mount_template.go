package d1pg

import (
	"context"
	"fmt"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) MountTemplates(ctx context.Context) (map[int]d1.MountTemplate, error) {
	return r.mountTemplates(ctx, "")
}

func (r *Repo) mountTemplates(ctx context.Context, conditions string, args ...interface{}) (map[int]d1.MountTemplate, error) {
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

	mountTemplates := make(map[int]d1.MountTemplate)
	for rows.Next() {
		var mountTemplate d1.MountTemplate
		var color1 *d1typ.Color
		var color2 *d1typ.Color
		var color3 *d1typ.Color
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

		maxEffects, err := d1.DecodeItemEffects(maxEffectsSli)
		if err != nil {
			return nil, err
		}
		mountTemplate.MaxEffects = maxEffects

		mountTemplates[mountTemplate.Id] = mountTemplate
	}

	return mountTemplates, nil
}
