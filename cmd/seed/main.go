package main

import (
	"context"
	"log"

	"github.com/boilerplate/internal/config"
	"github.com/boilerplate/internal/modules/permission/infra/persistence"
	"github.com/boilerplate/internal/platform/database"
	"github.com/boilerplate/pkg/logger"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(context.Background(), cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	seedData, err := LoadSeedData("configs/permissions.json")
	if err != nil {
		log.Fatalf("Failed to load seed data: %v", err)
	}

	seedLogger := logger.NewConsoleLogger()

	permRepo := persistence.NewPermissionRepoImpl(db)

	seedLogger.Info("seeding permissions...")
	if err := SeedPermissions(ctx, db, seedData); err != nil {
		seedLogger.Error("failed to seed permissions", "error", err)
		log.Fatalf("seed failed")
	}

	seedLogger.Info("permissions	 seeded successfully")

	// Check if permissions exist
	existingPerms, err := permRepo.FindAll(ctx)
	if err != nil {
		seedLogger.Error("failed to check existing permissions", "error", err)
		log.Fatalf("seed failed")
	}
	seedLogger.Info("existing permissions", "count", len(existingPerms))

	permissionIDs, err := LoadPermissionsByName(ctx, permRepo)
	if err != nil {
		seedLogger.Error("failed to load permissions", "error", err)
		log.Fatalf("seed failed")
	}

	seedLogger.Info("loaded permissions", "count", len(permissionIDs))

	seedLogger.Info("seeding group permissions", "count", len(seedData.GroupPermissions))
	groupIDs, err := SeedGroupPermissions(ctx, db, permissionIDs, seedData)
	if err != nil {
		seedLogger.Error("failed to seed group permissions", "error", err)
	}

	seedLogger.Info("seeding admin user...")
	if err := SeedAdminUser(ctx, db, seedData); err != nil {
		seedLogger.Error("failed to seed admin user", "error", err)
	}
	seedLogger.Info("admin user seeded successfully")

	seedLogger.Info("group permissions seeded successfully", "count", len(groupIDs))
}
