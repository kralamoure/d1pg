package retropg

import (
	"context"
	"fmt"
	"time"

	"github.com/kralamoure/retro"
	"github.com/kralamoure/retro/retrotyp"
)

func (r *Storer) CreateMount(ctx context.Context, mount retro.Mount) (id int, err error) {
	query := "INSERT INTO retro.mounts (template_id, character_id, name, sex, xp, capacities, validity)" +
		" VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		" RETURNING id;"

	var characterId *int
	if mount.CharacterId != 0 {
		characterId = &mount.CharacterId
	}

	var validity *time.Time
	if !mount.Validity.IsZero() {
		validity = &mount.Validity
	}

	capacities := make([]int, len(mount.Capacities))
	for i, v := range mount.Capacities {
		capacities[i] = int(v)
	}

	err = storerError(
		r.pool.QueryRow(ctx, query,
			mount.TemplateId, characterId, mount.Name, mount.Sex, mount.XP, capacities, validity,
		).Scan(&id),
	)
	return
}

func (r *Storer) UpdateMount(ctx context.Context, mount retro.Mount) error {
	query := "UPDATE retro.mounts" +
		" SET template_id = $2, character_id = $3, name = $4, sex = $5, xp = $6, capacities = $7, validity = $8" +
		" WHERE id = $1;"

	var characterId *int
	if mount.CharacterId != 0 {
		characterId = &mount.CharacterId
	}

	var validity *time.Time
	if !mount.Validity.IsZero() {
		validity = &mount.Validity
	}

	capacities := make([]int, len(mount.Capacities))
	for i, v := range mount.Capacities {
		capacities[i] = int(v)
	}

	tag, err := r.pool.Exec(ctx, query, mount.Id, mount.TemplateId, characterId, mount.Name, mount.Sex, mount.XP,
		capacities, validity)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return retro.ErrNotFound
	}

	return nil
}

func (r *Storer) DeleteMount(ctx context.Context, id int) error {
	query := "DELETE FROM retro.mounts" +
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

func (r *Storer) Mount(ctx context.Context, id int) (retro.Mount, error) {
	var mount retro.Mount

	mounts, err := r.mounts(ctx, "id = $1", id)
	if err != nil {
		return mount, err
	}

	if len(mounts) != 1 {
		return mount, retro.ErrNotFound
	}

	for k := range mounts {
		mount = mounts[k]
	}

	return mount, nil
}

func (r *Storer) Mounts(ctx context.Context) (items map[int]retro.Mount, err error) {
	return r.mounts(ctx, "")
}

func (r *Storer) MountsByCharacterId(ctx context.Context, characterId int) (items map[int]retro.Mount, err error) {
	return r.mounts(ctx, "character_id = $1", characterId)
}

func (r *Storer) mounts(ctx context.Context, conditions string, args ...interface{}) (map[int]retro.Mount, error) {
	query := "SELECT id, template_id, character_id, name, sex, xp, capacities, validity" +
		" FROM retro.mounts"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mounts := make(map[int]retro.Mount)
	for rows.Next() {
		var mount retro.Mount
		var characterId *int
		var validity *time.Time
		var capacities []int

		err = rows.Scan(&mount.Id, &mount.TemplateId, &characterId, &mount.Name, &mount.Sex, &mount.XP,
			&capacities, &validity)
		if err != nil {
			return nil, err
		}

		if characterId != nil {
			mount.CharacterId = *characterId
		}

		if validity != nil {
			mount.Validity = *validity
		}

		mount.Capacities = make([]retrotyp.MountCapacityId, len(capacities))
		for i, v := range capacities {
			mount.Capacities[i] = retrotyp.MountCapacityId(v)
		}

		mounts[mount.Id] = mount
	}

	return mounts, nil
}
