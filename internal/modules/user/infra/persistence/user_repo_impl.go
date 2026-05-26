package persistence

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/user/infra/persistence/dao"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var userAllowedFields = map[string]string{
	"id":         "id",
	"name":       "name",
	"email":      "email",
	"phone":      "phone",
	"role":       "role",
	"status":     "status",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var userColumns = []string{
	"id", "name", "email", "phone", "role", "password_hash", "token_version", "status", "created_at", "updated_at",
}

func scanUser(row pgx.Row) (*dao.UserDAO, error) {
	var user dao.UserDAO
	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.PasswordHash, &user.TokenVersion, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type UserRepoImpl struct {
	db *pgxpool.Pool
}

func NewUserRepoImpl(db *pgxpool.Pool) *UserRepoImpl {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) Save(ctx context.Context, user *domain.User) error {
	d := dao.FromDomain(user)
	query, args, err := psql.Insert("users").
		Columns("id", "name", "email", "phone", "role", "password_hash", "token_version", "status", "created_at", "updated_at").
		Values(d.ID, d.Name, d.Email, d.Phone, d.Role, d.PasswordHash, d.TokenVersion, d.Status, d.CreatedAt, d.UpdatedAt).
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

func (r *UserRepoImpl) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query, args, err := psql.Select(userColumns...).
		From("users").
		Where(sq.Eq{"id": uuid.UUID(id)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	userDAO, err := scanUser(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app_errors.NotFound("user")
		}
		return nil, app_errors.DatabaseFailure(err)
	}

	return userDAO.ToDomain(), nil
}

func (r *UserRepoImpl) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.User, int64, error) {
	// Count
	countBuilder := psql.Select("COUNT(*)").From("users")
	countBuilder, err := query.ApplyFilters(countBuilder, filters, userAllowedFields)
	if err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}

	// Fetch
	queryBuilder := psql.Select(userColumns...).From("users")
	queryBuilder, err = query.ApplyFilters(queryBuilder, filters, userAllowedFields)
	if err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}
	if len(sorts) > 0 {
		queryBuilder = query.ApplySort(queryBuilder, sorts, userAllowedFields)
	} else {
		queryBuilder = queryBuilder.OrderBy("created_at DESC")
	}

	dbQuery, args, err := queryBuilder.Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, dbQuery, args...)
	if err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		userDAO, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, userDAO.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}

	return users, total, nil
}

func (r *UserRepoImpl) Update(ctx context.Context, user *domain.User) error {
	d := dao.FromDomain(user)
	query, args, err := psql.Update("users").
		Set("name", d.Name).
		Set("email", d.Email).
		Set("phone", d.Phone).
		Set("role", d.Role).
		Set("status", d.Status).
		Set("token_version", d.TokenVersion).
		Set("updated_at", d.UpdatedAt).
		Where(sq.Eq{"id": d.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}

	if result.RowsAffected() == 0 {
		return app_errors.NotFound("user")
	}

	return nil
}

func (r *UserRepoImpl) Delete(ctx context.Context, id domain.UserID) error {
	query, args, err := psql.Delete("users").Where(sq.Eq{"id": uuid.UUID(id)}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}
	return nil
}

func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query, args, err := psql.Select(userColumns...).
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	userDAO, err := scanUser(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app_errors.NotFound("user")
		}
		return nil, app_errors.DatabaseFailure(err)
	}

	return userDAO.ToDomain(), nil
}
