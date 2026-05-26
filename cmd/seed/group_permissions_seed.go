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

func SeedGroupPermissions(ctx context.Context, db *pgxpool.Pool, permissionNames map[string]int, data *SeedData) (map[string]int, error) {
	if data == nil || len(data.GroupPermissions) == 0 {
		return nil, nil
	}

	groupPermRepo := persistence.NewGroupPermissionRepoImpl(db)
	// Check if permissions exist
	if len(permissionNames) == 0 {
		return nil, fmt.Errorf("no permissions found to associate with group permissions")
	}
	// Check existing group per	missions
	existing, err := groupPermRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing group permissions: %w", err)
	}

	fmt.Printf("Found %d existing group permissions\n", len(existing))
	for _, gp := range existing {
		fmt.Printf("Existing group: %s with %d permissions\n", gp.Name, len(gp.PermissionIDs))
	}

	for _, item := range data.GroupPermissions {
		permIDs := make([]int, 0, len(item.Permissions))
		for _, permissionName := range item.Permissions {
			id, ok := permissionNames[permissionName]
			if !ok {
				return nil, fmt.Errorf("permission %q not found for group %s", permissionName, item.Name)
			}
			permIDs = append(permIDs, id)
		}

		groupPerm := &domain.GroupPermission{
			Name:          item.Name,
			PermissionIDs: permIDs,
			CreatedAt:     time.Now(),
		}

		if err := groupPermRepo.Save(ctx, groupPerm); err != nil {
			return nil, fmt.Errorf("failed to save group permission %s: %w", item.Name, err)
		}
	}

	groups, err := groupPermRepo.FindAll(ctx)
	if err != nil {
		return nil, app_errors.DatabaseFailure(err)
	}

	result := make(map[string]int, len(groups))
	for _, gp := range groups {
		result[gp.Name] = gp.ID
	}
	return result, nil
}
