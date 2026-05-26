package persistence

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/permission/infra/persistence/dao"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PermissionRepoImpl struct {
	db *pgxpool.Pool
}

func NewPermissionRepoImpl(db *pgxpool.Pool) *PermissionRepoImpl {
	return &PermissionRepoImpl{db: db}
}

func (r *PermissionRepoImpl) FindAll(ctx context.Context) ([]*domain.Permission, error) {
	query, args, err := psql.Select(
		"id", "module", "action", "name", "description", "created_at",
	).From("permissions").OrderBy("id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, app_errors.DatabaseFailure(err)
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		var d dao.PermissionDAO
		err = rows.Scan(&d.ID, &d.Module, &d.Action, &d.Name, &d.Description, &d.CreatedAt)
		if err != nil {
			return nil, app_errors.DatabaseFailure(err)
		}
		permissions = append(permissions, d.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, app_errors.DatabaseFailure(err)
	}

	return permissions, nil
}

func (r *PermissionRepoImpl) SaveAll(ctx context.Context, permissions []*domain.Permission) error {
	if len(permissions) == 0 {
		return nil
	}

	builder := psql.Insert("permissions").
		Columns("module", "action", "name", "description")

	for _, p := range permissions {
		builder = builder.Values(p.Module, p.Action, p.Name, p.Description)
	}

	query, args, err := builder.
		Suffix("ON CONFLICT (name) DO UPDATE SET module = EXCLUDED.module, action = EXCLUDED.action, description = EXCLUDED.description").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}
	return nil
}

type GroupPermissionRepoImpl struct {
	db *pgxpool.Pool
}

func NewGroupPermissionRepoImpl(db *pgxpool.Pool) *GroupPermissionRepoImpl {
	return &GroupPermissionRepoImpl{db: db}
}

func (r *GroupPermissionRepoImpl) FindByName(ctx context.Context, name string) (*domain.GroupPermission, error) {
	query, args, err := psql.Select("id", "name", "permission_ids", "created_at").
		From("group_permissions").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var d dao.GroupPermissionDAO
	err = r.db.QueryRow(ctx, query, args...).Scan(&d.ID, &d.Name, pq.Array(&d.PermissionIDs), &d.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app_errors.NotFound("group_permission")
		}
		return nil, app_errors.DatabaseFailure(err)
	}

	return d.ToDomain(), nil
}

func (r *GroupPermissionRepoImpl) FindAll(ctx context.Context) ([]*domain.GroupPermission, error) {
	query, args, err := psql.Select("id", "name", "permission_ids", "created_at").
		From("group_permissions").
		OrderBy("name").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, app_errors.DatabaseFailure(err)
	}
	defer rows.Close()

	var groups []*domain.GroupPermission
	for rows.Next() {
		var d dao.GroupPermissionDAO
		err = rows.Scan(&d.ID, &d.Name, pq.Array(&d.PermissionIDs), &d.CreatedAt)
		if err != nil {
			return nil, app_errors.DatabaseFailure(err)
		}
		groups = append(groups, d.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, app_errors.DatabaseFailure(err)
	}

	return groups, nil
}

func (r *GroupPermissionRepoImpl) Save(ctx context.Context, gp *domain.GroupPermission) error {
	d := dao.FromGroupDomain(gp)
	query, args, err := psql.Insert("group_permissions").
		Columns("name", "permission_ids").
		Values(d.Name, pq.Array(d.PermissionIDs)).
		Suffix("ON CONFLICT (name) DO UPDATE SET permission_ids = EXCLUDED.permission_ids").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}
	return nil
}
