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
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'booking_status') THEN
				CREATE TYPE booking_status AS ENUM ('pending', 'approved', 'rejected');
			END IF;
		END
		$$;
	`).Error
	if err != nil {
		return err
	}

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

	createEventDetailsView := `
	CREATE OR REPLACE VIEW event_details AS
		SELECT
			e.id,
			e.name,
			e.description,
			e.start_time,
			e.end_time,
			e.event_type,
			e.created_by AS creator_id,
			u.name AS creator_name,
			e.created_at,
			e.updated_at,
			e.deleted_at,
			e.created_by
		FROM
			events e
		LEFT JOIN
			users u ON e.created_by = u.id;
	`
	if err := db.Exec(createEventDetailsView).Error; err != nil {
		return err
	}

	fullInvitationDetailsView := `
	CREATE OR REPLACE VIEW full_invitation_details AS
	SELECT
		i.id AS id, 
		e.id AS event_id,
		e.name AS event_name,
		u.id AS user_id,
		u.name AS user_name,
		u.email AS user_email,
		ui.invited_at,
		ui.rsvp_status,
		ui.rsvp_at,
		ui.attended_at,
		ui.qr_code
	FROM
		invitations i
	JOIN
		user_invitation ui ON i.id = ui.invitation_id
	JOIN
		users u ON ui.user_id = u.id
	JOIN
		events e ON i.event_id = e.id;
	`
	if err := db.Exec(fullInvitationDetailsView).Error; err != nil {
		return err
	}

	ormawaEventsView := `
	-- Events created by "Ormawa" View
	CREATE OR REPLACE VIEW ormawa_events_view AS
	SELECT
		e.id AS event_id,
		e.name AS event_name,
		e.description,
		e.start_time,
		e.end_time,
		u.id AS creator_id,
		u.name AS creator_name
	FROM
		events e
	JOIN
		users u ON e.created_by = u.id
	WHERE
		u.role = 'ormawa' AND e.deleted_at IS NULL;
	`
	if err := db.Exec(ormawaEventsView).Error; err != nil {
		return err
	}

	if err := db.Exec(ormawaEventsView).Error; err != nil {
		return err
	}

	userAttendanceView := `
	-- User Attendance View
	CREATE OR REPLACE VIEW user_attendance_view AS
	SELECT
		u.id as user_id,
		u.name as user_name,
		e.id as event_id,
		e.name as event_name,
		ui.attended_at
	FROM
		users u
	JOIN
		user_invitation ui ON u.id = ui.user_id
	JOIN
		invitations i ON ui.invitation_id = i.id
	JOIN
		events e ON i.event_id = e.id
	WHERE 
		ui.attended_at IS NOT NULL
	AND
		u.deleted_at IS NULL;
	`
	if err := db.Exec(userAttendanceView).Error; err != nil {
		return err
	}
	return nil
}
