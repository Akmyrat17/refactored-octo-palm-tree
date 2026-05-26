package persistance

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/profile/infra/persistance/dao"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var profileAllowedFields = map[string]string{
	"id":         "id",
	"full_name":  "full_name",
	"avatar_url": "avatar_url",
	"headline":   "headline",
	"bio":        "bio",
	"resume_url": "resume_url",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var profileColumns = []string{
	"id", "full_name", "avatar_url", "headline", "bio", "resume_url", "social_links", "created_at", "updated_at",
}

func scanProfile(row pgx.Row) (*dao.ProfileDAO, error) {
	var profile dao.ProfileDAO
	var socialLinksJSON []byte

	err := row.Scan(
		&profile.ID, &profile.FullName, &profile.AvatarURL, &profile.Headline, &profile.Bio,
		&profile.ResumeURL, &socialLinksJSON, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(socialLinksJSON) > 0 {
		if err := json.Unmarshal(socialLinksJSON, &profile.SocialLinks); err != nil {
			return nil, err
		}
	}

	return &profile, nil
}

type ProfileRepoImpl struct {
	db *pgxpool.Pool
}

func NewProfileRepoImpl(db *pgxpool.Pool) *ProfileRepoImpl {
	return &ProfileRepoImpl{db: db}
}

func (r *ProfileRepoImpl) Create(ctx context.Context, profile *dao.ProfileDAO) error {
	socialLinksJSON, err := json.Marshal(profile.SocialLinks)
	if err != nil {
		return err
	}

	query, args, err := psql.Insert("profiles").
		Columns("id", "full_name", "avatar_url", "headline", "bio", "resume_url", "social_links", "created_at", "updated_at").
		Values(profile.ID, profile.FullName, profile.AvatarURL, profile.Headline, profile.Bio, profile.ResumeURL, socialLinksJSON, profile.CreatedAt, profile.UpdatedAt).
		ToSql()

	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfileRepoImpl) FindByID(ctx context.Context, id domain.ProfileID) (*domain.Profile, error) {
	query, args, err := psql.Select(profileColumns...).
		From("profiles").
		Where(sq.Eq{"id": uuid.UUID(id)}).
		ToSql()
	if err != nil {
		return nil, err
	}

	profileDAO, err := scanProfile(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app_errors.NotFound("profile")
		}
		return nil, app_errors.DatabaseFailure(err)
	}

	return profileDAO.ToDomain(), nil
}

func (r *ProfileRepoImpl) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Profile, int64, error) {
	// Count
	countBuilder := psql.Select("COUNT(*)").From("profiles")
	countBuilder, err := query.ApplyFilters(countBuilder, filters, profileAllowedFields)
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
	queryBuilder := psql.Select(profileColumns...).From("profiles")
	queryBuilder, err = query.ApplyFilters(queryBuilder, filters, profileAllowedFields)
	if err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}
	if len(sorts) > 0 {
		queryBuilder = query.ApplySort(queryBuilder, sorts, profileAllowedFields)
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

	var profiles []*domain.Profile

	for rows.Next() {
		profileDAO, err := scanProfile(rows)
		if err != nil {
			return nil, 0, err
		}
		profiles = append(profiles, profileDAO.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}

	return profiles, total, nil
}

func (r *ProfileRepoImpl) Update(ctx context.Context, profile *dao.ProfileDAO) error {
	socialLinksJSON, err := json.Marshal(profile.SocialLinks)
	if err != nil {
		return err
	}

	setMap := sq.Eq{
		"full_name":    sq.Expr("COALESCE(NULLIF(?::text, ''), full_name)", profile.FullName),
		"avatar_url":   sq.Expr("COALESCE(NULLIF(?::text, ''), avatar_url)", profile.AvatarURL),
		"headline":     sq.Expr("COALESCE(NULLIF(?::text, ''), headline)", profile.Headline),
		"bio":          sq.Expr("COALESCE(NULLIF(?::text, ''), bio)", profile.Bio),
		"resume_url":   sq.Expr("COALESCE(NULLIF(?::text, ''), resume_url)", profile.ResumeURL),
		"social_links": socialLinksJSON,
		"updated_at":   time.Now(),
	}

	query, args, err := psql.Update("profiles").
		SetMap(setMap).
		Where(sq.Eq{"id": profile.ID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}

	if result.RowsAffected() == 0 {
		return app_errors.NotFound("profile")
	}

	return nil
}

func (r *ProfileRepoImpl) Delete(ctx context.Context, id domain.ProfileID) error {
	query, args, err := psql.Delete("profiles").
		Where(sq.Eq{"id": uuid.UUID(id)}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}

	if result.RowsAffected() == 0 {
		return app_errors.NotFound("profile")
	}

	return nil
}
