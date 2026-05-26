package main

import (
	"context"
	"fmt"
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/permission/infra/persistence"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedPermissions(ctx context.Context, db *pgxpool.Pool, data *SeedData) error {
	if data == nil || len(data.Permissions) == 0 {
		return nil
	}
	permRepo := persistence.NewPermissionRepoImpl(db)
	// Check existing permissions
	existing, err := permRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing permissions: %w", err)
	}

	existingNames := make(map[string]bool)
	for _, p := range existing {
		existingNames[p.Name] = true
	}

	// Filter out existing permissions
	var newPermissions []*domain.Permission
	for _, item := range data.Permissions {
		if !existingNames[item.Name] {
			newPermissions = append(newPermissions, &domain.Permission{
				Module:      item.Module,
				Action:      item.Action,
				Name:        item.Name,
				Description: item.Description,
				CreatedAt:   time.Now(),
			})
		}
	}

	if len(newPermissions) == 0 {
		fmt.Println("All permissions already exist, skipping...")
		return nil
	}

	fmt.Printf("Seeding %d new permissions...\n", len(newPermissions))
	if err := permRepo.SaveAll(ctx, newPermissions); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	return nil
}

func LoadPermissionsByName(ctx context.Context, db *persistence.PermissionRepoImpl) (map[string]int, error) {
	permissions, err := db.FindAll(ctx)
	if err != nil {
		return nil, app_errors.DatabaseFailure(err)
	}

	result := make(map[string]int, len(permissions))
	for _, p := range permissions {
		result[p.Name] = p.ID
	}
	return result, nil
}
