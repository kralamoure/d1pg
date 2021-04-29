package d1pg

import (
	"context"
	"fmt"
	"time"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) CreateMount(ctx context.Context, mount d1.Mount) (id int, err error) {
	query := "INSERT INTO d1.mounts (template_id, character_id, name, sex, xp, capacities, validity)" +
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

	err = repoError(
		r.pool.QueryRow(ctx, query,
			mount.TemplateId, characterId, mount.Name, mount.Sex, mount.XP, capacities, validity,
		).Scan(&id),
	)
	return
}

func (r *Repo) UpdateMount(ctx context.Context, mount d1.Mount) error {
	query := "UPDATE d1.mounts" +
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
		return d1.ErrNotFound
	}

	return nil
}

func (r *Repo) DeleteMount(ctx context.Context, id int) error {
	query := "DELETE FROM d1.mounts" +
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

func (r *Repo) Mount(ctx context.Context, id int) (d1.Mount, error) {
	var mount d1.Mount

	mounts, err := r.mounts(ctx, "id = $1", id)
	if err != nil {
		return mount, err
	}

	if len(mounts) != 1 {
		return mount, d1.ErrNotFound
	}

	for k := range mounts {
		mount = mounts[k]
	}

	return mount, nil
}

func (r *Repo) Mounts(ctx context.Context) (items map[int]d1.Mount, err error) {
	return r.mounts(ctx, "")
}

func (r *Repo) MountsByCharacterId(ctx context.Context, characterId int) (items map[int]d1.Mount, err error) {
	return r.mounts(ctx, "character_id = $1", characterId)
}

func (r *Repo) mounts(ctx context.Context, conditions string, args ...interface{}) (map[int]d1.Mount, error) {
	query := "SELECT id, template_id, character_id, name, sex, xp, capacities, validity" +
		" FROM d1.mounts"
	if conditions != "" {
		query += fmt.Sprintf(" WHERE %s", conditions)
	}
	query += ";"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mounts := make(map[int]d1.Mount)
	for rows.Next() {
		var mount d1.Mount
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

		mount.Capacities = make([]d1typ.MountCapacityId, len(capacities))
		for i, v := range capacities {
			mount.Capacities[i] = d1typ.MountCapacityId(v)
		}

		mounts[mount.Id] = mount
	}

	return mounts, nil
}
