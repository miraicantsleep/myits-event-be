package migrations

import (
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Ensure enum type exists before AutoMigrate
	err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
				CREATE TYPE user_role AS ENUM ('user', 'departemen', 'ormawa', 'admin');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
				CREATE TYPE event_type AS ENUM ('online', 'offline');
			END IF;
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'rsvp_status') THEN
				CREATE TYPE rsvp_status AS ENUM ('accepted', 'declined', 'pending');
			END IF;
		END
		$$;
	`).Error
	if err != nil {
		return err
	}

	if err := db.SetupJoinTable(&entity.Invitation{}, "Users", &entity.UserInvitation{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.User{}, &entity.Department{}, &entity.Event{}, &entity.Room{}, &entity.Invitation{}, &entity.BookingRequest{},
	); err != nil {
		return err
	}

	return nil
}
