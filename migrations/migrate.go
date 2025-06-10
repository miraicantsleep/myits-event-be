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

	// SQL for QR Code generation function and trigger
	qrCodeFunctionSQL := `
	CREATE OR REPLACE FUNCTION generate_user_invitation_qr_code()
	RETURNS TRIGGER AS $$
	BEGIN
	    -- Check if qr_code is NULL or an empty string, then generate.
	    -- This allows for manually setting a QR code if ever needed,
	    -- though typically it will be NULL on insert.
	    IF NEW.qr_code IS NULL OR NEW.qr_code = '' THEN
	        NEW.qr_code := uuid_generate_v4();
	    END IF;
	    RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(qrCodeFunctionSQL).Error; err != nil {
		return err
	}

	// Setup Join Table and AutoMigrate
	if err := db.SetupJoinTable(&entity.Invitation{}, "Users", &entity.UserInvitation{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&entity.User{}, &entity.Department{}, &entity.Event{}, &entity.Room{}, &entity.Invitation{}, &entity.BookingRequest{}, &entity.UserInvitation{},
	); err != nil {
		return err
	}

	// SQL for QR Code trigger - MOVED TO AFTER AUTOMIGRATE
	qrCodeTriggerSQL := `
	DROP TRIGGER IF EXISTS trg_generate_qr_code_before_insert_on_user_invitation ON user_invitation;
	CREATE TRIGGER trg_generate_qr_code_before_insert_on_user_invitation
	    BEFORE INSERT ON user_invitation
	    FOR EACH ROW
	    EXECUTE FUNCTION generate_user_invitation_qr_code();
	`
	if err := db.Exec(qrCodeTriggerSQL).Error; err != nil {
		return err
	}

	return nil
}
