package main

import (
	"context"
	"fmt"
	"time"

	"github.com/boilerplate/internal/shared/enums"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser(ctx context.Context, db *pgxpool.Pool, data *SeedData) error {
	if data == nil || data.AdminUser.Email == "" {
		return nil
	}

	if data.AdminUser.Password == "" {
		return fmt.Errorf("admin password is required in seed data")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.AdminUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	groupPermissionID, err := loadGroupPermissionID(ctx, db, data.AdminUser.GroupPermission)
	if err != nil {
		return err
	}

	if groupPermissionID == nil {
		return fmt.Errorf("admin group permission %q not found", data.AdminUser.GroupPermission)
	}

	query := `INSERT INTO users
	(id, name, email, phone, role, password_hash, token_version, status, group_permission_id, created_at, updated_at)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (email) DO UPDATE SET
	name = EXCLUDED.name,
	phone = EXCLUDED.phone,
	role = EXCLUDED.role,
	password_hash = EXCLUDED.password_hash,
	token_version = EXCLUDED.token_version,
	status = EXCLUDED.status,
	group_permission_id = EXCLUDED.group_permission_id,
	updated_at = EXCLUDED.updated_at;`

	id := uuid.New().String()
	_, err = db.Exec(ctx, query,
		id,
		data.AdminUser.Name,
		data.AdminUser.Email,
		data.AdminUser.Phone,
		enums.RoleAdmin.String(),
		hash,
		1,
		"active",
		*groupPermissionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	return nil
}

func loadGroupPermissionID(ctx context.Context, db *pgxpool.Pool, name string) (*int, error) {
	var id int
	query := `SELECT id FROM group_permissions WHERE name = $1 LIMIT 1`
	if err := db.QueryRow(ctx, query, name).Scan(&id); err != nil {
		return nil, fmt.Errorf("failed to find group permission %q: %w", name, err)
	}
	return &id, nil
}
