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

	// Trigger to prevent modification of approved/rejected bookings
	preventBookingModificationFunc := `
	CREATE OR REPLACE FUNCTION prevent_booking_modification()
	RETURNS TRIGGER AS $$
	BEGIN
	    IF OLD.status IN ('approved', 'rejected') THEN
	        RAISE EXCEPTION 'Cannot modify a booking request that has already been approved or rejected.';
	    END IF;
	    RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(preventBookingModificationFunc).Error; err != nil {
		return err
	}

	preventBookingModificationTrigger := `
	DROP TRIGGER IF EXISTS trg_prevent_booking_modification ON booking_requests;
	CREATE TRIGGER trg_prevent_booking_modification
	BEFORE UPDATE ON booking_requests
	FOR EACH ROW EXECUTE FUNCTION prevent_booking_modification();
	`
	if err := db.Exec(preventBookingModificationTrigger).Error; err != nil {
		return err
	}

	// Trigger to prevent deletion of events with active bookings
	preventEventDeletionFunc := `
	CREATE OR REPLACE FUNCTION prevent_event_deletion_with_bookings()
	RETURNS TRIGGER AS $$
	BEGIN
	    IF EXISTS (SELECT 1 FROM booking_requests WHERE event_id = OLD.id AND status IN ('pending', 'approved') AND deleted_at IS NULL) THEN
	        RAISE EXCEPTION 'Cannot delete event: It has active booking requests.';
	    END IF;
	    RETURN OLD;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(preventEventDeletionFunc).Error; err != nil {
		return err
	}

	preventEventDeletionTrigger := `
	DROP TRIGGER IF EXISTS trg_prevent_event_deletion_with_bookings ON events;
	CREATE TRIGGER trg_prevent_event_deletion_with_bookings
	BEFORE DELETE ON events
	FOR EACH ROW EXECUTE FUNCTION prevent_event_deletion_with_bookings();
	`
	if err := db.Exec(preventEventDeletionTrigger).Error; err != nil {
		return err
	}

	getEventAttendeesFunc := `
	CREATE OR REPLACE FUNCTION get_event_attendees(p_event_id uuid)
	RETURNS TABLE (
		user_id uuid,
		user_name character varying,
		user_email character varying,
		attended_at timestamp
	) AS $$
	BEGIN
		RETURN QUERY
		SELECT
			u.id,
			u.name,
			u.email,
			ui.attended_at
		FROM
			user_invitation ui
		JOIN
			users u ON ui.user_id = u.id
		JOIN
			invitations i ON ui.invitation_id = i.id
		WHERE
			i.event_id = p_event_id
			AND ui.attended_at IS NOT NULL
			AND u.deleted_at IS NULL;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(getEventAttendeesFunc).Error; err != nil {
		return err
	}

	getPendingBookingsFunc := `
	CREATE OR REPLACE FUNCTION get_pending_booking_requests_for_department(p_department_id uuid)
	RETURNS TABLE (
		booking_request_id uuid,
		event_name character varying,
		event_start_time timestamp,
		event_end_time timestamp,
		room_name character varying,
		requesting_ormawa character varying
	) AS $$
	BEGIN
		RETURN QUERY
		SELECT
			br.id,
			e.name,
			e.start_time,
			e.end_time,
			r.name,
			u.name
		FROM
			booking_requests br
		JOIN
			events e ON br.event_id = e.id
		JOIN
			users u ON e.created_by = u.id
		JOIN
			booking_request_room brr ON br.id = brr.booking_request_id
		JOIN
			rooms r ON brr.room_id = r.id
		WHERE
			r.department_id = p_department_id
			AND br.status = 'pending'
			AND br.deleted_at IS NULL;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(getPendingBookingsFunc).Error; err != nil {
		return err
	}

	createBookingWithRoomsView := `
	CREATE OR REPLACE VIEW vw_booking_with_rooms AS
	SELECT
		br.id AS booking_id,
		br.status AS booking_status,
		e.id AS event_id,
		e.name AS event_name,
		r.id AS room_id,
		r.name AS room_name,
		u.name AS requested_by
	FROM
		booking_requests br
	JOIN
		events e ON br.event_id = e.id
	JOIN
		users u ON e.created_by = u.id
	JOIN
		booking_request_room brr ON br.id = brr.booking_request_id
	JOIN
		rooms r ON brr.room_id = r.id;
	`
	if err := db.Exec(createBookingWithRoomsView).Error; err != nil {
		return err
	}

	roomDetailsView := `
	CREATE OR REPLACE VIEW vw_room_details AS
	SELECT
		r.id,
		r.name,
		r.capacity,
		r.department_id,
		d.name AS department_name,
		r.created_at,
		r.updated_at,
		r.deleted_at
	FROM
		rooms r
	LEFT JOIN
		departments d ON r.department_id = d.id
	WHERE
		r.deleted_at IS NULL;
	`
	if err := db.Exec(roomDetailsView).Error; err != nil {
		return err
	}

	setInvitedAtFunction := `
	CREATE OR REPLACE FUNCTION fn_set_invited_at_timestamp()
	RETURNS TRIGGER AS $$
	BEGIN
		-- Set kolom invited_at dengan waktu transaksi saat ini
		NEW.invited_at := NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(setInvitedAtFunction).Error; err != nil {
		return err
	}

	setInvitedAtTrigger := `
	DROP TRIGGER IF EXISTS trg_auto_set_invited_at ON user_invitation;
	CREATE TRIGGER trg_auto_set_invited_at
	BEFORE INSERT ON user_invitation
	FOR EACH ROW
	EXECUTE FUNCTION fn_set_invited_at_timestamp();
	`

	if err := db.Exec(setInvitedAtTrigger).Error; err != nil {
		return err
	}

	getEventByStatusFunc := `
	CREATE OR REPLACE FUNCTION get_event_by_status(p_timeline_status TEXT)
	RETURNS TABLE (
		id uuid,
		name character varying,
		description text,
		start_time timestamp,
		end_time timestamp,
		event_type event_type,
		creator_name character varying
	) AS $$
	BEGIN
		IF p_timeline_status = 'ongoing' THEN
			RETURN QUERY
			SELECT e.id, e.name, e.description, e.start_time, e.end_time, e.event_type, u.name
			FROM events e
			JOIN users u ON e.created_by = u.id
			WHERE e.deleted_at IS NULL AND NOW() BETWEEN e.start_time AND e.end_time;


		ELSIF p_timeline_status = 'upcoming' THEN
			RETURN QUERY
			SELECT e.id, e.name, e.description, e.start_time, e.end_time, e.event_type, u.name
			FROM events e
			JOIN users u ON e.created_by = u.id
			WHERE e.deleted_at IS NULL AND e.start_time > NOW();


		ELSIF p_timeline_status = 'finished' THEN
			RETURN QUERY
			SELECT e.id, e.name, e.description, e.start_time, e.end_time, e.event_type, u.name
			FROM events e
			JOIN users u ON e.created_by = u.id
			WHERE e.deleted_at IS NULL AND e.end_time < NOW();


		ELSE
			RAISE EXCEPTION 'Invalid timeline status';
		END IF;
	END;
	$$ LANGUAGE plpgsql;
	`

	if err := db.Exec(getEventByStatusFunc).Error; err != nil {
		return err
	}

	isRoomAvailableFunc := `
	CREATE OR REPLACE FUNCTION is_room_available(
		p_room_id UUID,
		p_start_time TIMESTAMP,
		p_end_time TIMESTAMP
	)
	RETURNS BOOLEAN AS $$
	DECLARE
		is_available BOOLEAN;
	BEGIN
		SELECT NOT EXISTS (
			SELECT 1
			FROM booking_requests br
			JOIN booking_request_room brr ON br.id = brr.booking_request_id
			JOIN events e ON br.event_id = e.id
			WHERE brr.room_id = p_room_id
				AND br.status = 'approved'
				AND br.deleted_at IS NULL
				AND (p_start_time < e.end_time AND p_end_time > e.start_time)
		) INTO is_available;


		RETURN is_available;
	END;
	$$ LANGUAGE plpgsql;
	`

	if err := db.Exec(isRoomAvailableFunc).Error; err != nil {
		return err
	}

	validateEventTimeFunc := `
	CREATE OR REPLACE FUNCTION validate_event_time()
	RETURNS TRIGGER AS $$
	BEGIN
		IF NEW.end_time <= NEW.start_time THEN
			RAISE EXCEPTION 'end time harus lebih besar dari start time';
		END IF;
		-- Jika valid, lanjutkan operasi INSERT atau UPDATE
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(validateEventTimeFunc).Error; err != nil {
		return err
	}

	validateEventTimeTrigger := `
	DROP TRIGGER IF EXISTS trg_validate_event_time ON events;
	CREATE TRIGGER trg_validate_event_time
	BEFORE INSERT OR UPDATE ON events
	FOR EACH ROW
	EXECUTE FUNCTION validate_event_time();
	`
	if err := db.Exec(validateEventTimeTrigger).Error; err != nil {
		return err
	}

	calculateEventDurationFunc := `
	CREATE OR REPLACE FUNCTION calculate_event_duration()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.duration_in_minutes := EXTRACT(EPOCH FROM (NEW.end_time - NEW.start_time)) / 60;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(calculateEventDurationFunc).Error; err != nil {
		return err
	}
	calculateEventDurationTrigger := `
	DROP TRIGGER IF EXISTS trg_set_event_duration ON events;
	CREATE TRIGGER trg_set_event_duration
	BEFORE INSERT OR UPDATE ON events
	FOR EACH ROW
	EXECUTE FUNCTION calculate_event_duration();
	`
	if err := db.Exec(calculateEventDurationTrigger).Error; err != nil {
		return err
	}

	getUserUpcomingEventsFunc := `
	CREATE OR REPLACE FUNCTION get_user_upcoming_events(p_user_id UUID)
	RETURNS TABLE (
		event_id uuid,
		event_name character varying,
		event_description text,
		event_start_time timestamp,
		event_end_time timestamp,
		rsvp_status rsvp_status
	) AS $$
	BEGIN
		RETURN QUERY
		SELECT
			e.id,
			e.name,
			e.description,
			e.start_time,
			e.end_time,
			ui.rsvp_status
		FROM
			events e
		JOIN
			invitations i ON e.id = i.event_id
		JOIN
			user_invitation ui ON i.id = ui.invitation_id
		WHERE
			ui.user_id = p_user_id
			AND e.start_time > NOW()
			AND e.deleted_at IS NULL
		ORDER BY
			e.start_time ASC;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(getUserUpcomingEventsFunc).Error; err != nil {
		return err
	}

	getCreatedEventCountFunc := `
	CREATE OR REPLACE FUNCTION get_created_event_count(p_user_id UUID)
	RETURNS INT AS $$
	DECLARE
		event_count INT;
	BEGIN
		SELECT
			COUNT(*)
		INTO
			event_count
		FROM
			events
		WHERE
			created_by = p_user_id
			AND deleted_at IS NULL;
			
		RETURN event_count;
	END;
	$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(getCreatedEventCountFunc).Error; err != nil {
		return err
	}

	return nil
}
