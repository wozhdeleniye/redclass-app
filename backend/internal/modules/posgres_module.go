package modules

import (
	"context"
	"fmt"

	"github.com/wozhdeleniye/redclass-app/internal/config"
	"github.com/wozhdeleniye/redclass-app/internal/migrations"
	pgrepo "github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(conf config.DatabaseConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func RunGormMigrations(lc fx.Lifecycle, migrator *migrations.GormMigrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return migrator.Migrate()
		},
	})
}

var GormModule = fx.Module("gorm",
	fx.Provide(
		NewGormDB,
		migrations.NewGormMigrator,
		pgrepo.NewUserRepository,
	),
	fx.Invoke(RunGormMigrations),
)
