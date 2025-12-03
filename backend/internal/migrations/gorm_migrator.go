package migrations

import (
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type GormMigrator struct {
	db *gorm.DB
}

func NewGormMigrator(db *gorm.DB) *GormMigrator {
	return &GormMigrator{db: db}
}

func (m *GormMigrator) Migrate() error {
	m.db.Exec("SET CONSTRAINTS ALL DEFERRED")

	tables, err := m.db.Migrator().GetTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if err := m.db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	err = m.db.AutoMigrate(
		&models.User{},
		&models.Subject{},
		&models.Role{},
		&models.Task{},
		&models.Project{},
		&models.ProjectMember{},
	)
	return err
}
