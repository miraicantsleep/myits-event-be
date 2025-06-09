package migrations

import (
	"github.com/miraicantsleep/myits-event-be/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}

	if err := seeds.DepartmentAndUserSeeder(db); err != nil {
		return err
	}

	return nil
}
