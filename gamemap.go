package d1pg

import (
	"context"

	"github.com/kralamoure/d1"
)

func (r *Repo) GameMaps(ctx context.Context) (gameMaps map[int]d1.GameMap, err error) {
	query := "SELECT id, name, width, height, background, ambiance, music, outdoor, capabilities, data, encrypted_data, key" +
		" FROM d1_static.maps;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	gameMaps = make(map[int]d1.GameMap)
	for rows.Next() {
		var gameMap d1.GameMap

		err = rows.Scan(&gameMap.Id, &gameMap.Name, &gameMap.Width, &gameMap.Height,
			&gameMap.Background, &gameMap.Ambiance, &gameMap.Music, &gameMap.Outdoor, &gameMap.Capabilities,
			&gameMap.Data, &gameMap.EncryptedData, &gameMap.Key)
		if err != nil {
			return
		}

		gameMaps[gameMap.Id] = gameMap
	}
	return
}
