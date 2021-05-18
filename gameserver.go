package retropg

import (
	"context"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Db) CreateGameServer(ctx context.Context, gameServer retro.GameServer) error {
	query := "INSERT INTO retro.gameservers (id, host, port, state, completion)" +
		" VALUES ($1, $2, $3, $4, $5);"

	_, err := r.pool.Exec(ctx, query,
		gameServer.Id, gameServer.Host, gameServer.Port, gameServer.State, gameServer.Completion)
	return dbError(err)
}

func (r *Db) GameServers(ctx context.Context) (gameServers map[int]retro.GameServer, err error) {
	query := "SELECT id, host, port, state, completion" +
		" FROM retro.gameservers;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	gameServers = make(map[int]retro.GameServer)
	for rows.Next() {
		var gameServer retro.GameServer
		err = rows.Scan(&gameServer.Id, &gameServer.Host, &gameServer.Port, &gameServer.State, &gameServer.Completion)
		if err != nil {
			return
		}
		gameServers[gameServer.Id] = gameServer
	}
	return
}

func (r *Db) GameServer(ctx context.Context, id int) (gameServer retro.GameServer, err error) {
	query := "SELECT id, host, port, state, completion" +
		" FROM retro.gameservers" +
		" WHERE id = $1;"

	err = dbError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&gameServer.Id, &gameServer.Host, &gameServer.Port, &gameServer.State, &gameServer.Completion),
	)
	return
}

func (r *Db) SetGameServerState(ctx context.Context, id int, state retrotyp.GameServerState) error {
	query := "UPDATE retro.gameservers" +
		" SET state = $2" +
		" WHERE id = $1;"

	tag, err := r.pool.Exec(ctx, query, id, state)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return retro.ErrNotFound
	}
	return nil
}
