package d1pg

import (
	"context"

	"github.com/kralamoure/d1"
)

func (r *Repo) NPCs(ctx context.Context, gameServerId int) (npcs map[string]d1.NPC, err error) {
	query := "SELECT id, gameserver_id, map_id, cell_id, direction, template_id, sex, gfx, scale_x, scale_y, color_1, color_2, color_3, accessories, extra_clip, custom_artwork, dialog_id, market_id" +
		" FROM d1.npcs" +
		" WHERE gameserver_id = $1;"

	rows, err := r.pool.Query(ctx, query, gameServerId)
	if err != nil {
		return
	}
	defer rows.Close()

	npcs = make(map[string]d1.NPC)
	for rows.Next() {
		var npc d1.NPC
		var dialogId *int
		var marketId *string
		var extraClip *int
		var customArtwork *int

		err = rows.Scan(&npc.Id, &npc.GameServerId, &npc.MapId, &npc.CellId, &npc.Direction, &npc.TemplateId,
			&npc.Sex, &npc.GFX, &npc.ScaleX, &npc.ScaleY, &npc.Color1, &npc.Color2, &npc.Color3, &npc.Accessories,
			&extraClip, &customArtwork, &dialogId, &marketId)
		if err != nil {
			return
		}

		if dialogId != nil {
			npc.DialogId = *dialogId
		}

		if marketId != nil {
			npc.MarketId = *marketId
		}

		if extraClip != nil {
			npc.ExtraClip = *extraClip
		} else {
			npc.ExtraClip = -1
		}

		if customArtwork != nil {
			npc.CustomArtwork = *customArtwork
		}

		npcs[npc.Id] = npc
	}
	return
}
