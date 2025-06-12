package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/config" // Added for EmailConfig
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/utils"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type (
	InvitationService interface {
		Create(ctx context.Context, req dto.CreateInvitationRequest) (dto.CreateInvitationResponse, error)
		GetInvitationByID(ctx context.Context, invitationID string) ([]dto.InvitationResponse, error)
		GetInvitationByEventID(ctx context.Context, eventID string) ([]dto.InvitationResponse, error)
		GetInvitationByUserID(ctx context.Context, userID string) ([]dto.InvitationResponse, error)
		GetAllInvitations(ctx context.Context) ([]dto.InvitationResponse, error)
		Update(ctx context.Context, invitationID string, req dto.UpdateInvitationRequest) (dto.InvitationResponse, error)
		Delete(ctx context.Context, invitationID string) error
		ScanQRCode(ctx context.Context, qrCode string) (dto.ScanQRCodeResponse, error)
		ProcessRSVP(ctx context.Context, qrCodeToken string, newRsvpStatus string) error // New method
	}

	invitationService struct {
		invitationRepo repository.InvitationRepository
		jwtService     JWTService
		db             *gorm.DB
	}
)

func NewInvitationService(
	invitationRepo repository.InvitationRepository,
	jwtService JWTService,
	db *gorm.DB,
) InvitationService {
	return &invitationService{
		invitationRepo: invitationRepo,
		jwtService:     jwtService,
		db:             db,
	}
}

// Create handles invitation creation, skipping already-invited users
func (s *invitationService) Create(ctx context.Context, req dto.CreateInvitationRequest) (dto.CreateInvitationResponse, error) {
	// parse event ID
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return dto.CreateInvitationResponse{}, err
	}

	// prepare user entities
	users := make([]entity.User, len(req.UserIDs))
	for i, id := range req.UserIDs {
		uid, err := uuid.Parse(id)
		if err != nil {
			return dto.CreateInvitationResponse{}, err
		}
		users[i] = entity.User{ID: uid}
	}

	// filter out already-invited users
	var toInvite []entity.User
	for _, u := range users {
		exists, err := s.invitationRepo.CheckInvitationExist(ctx, nil, eventID, u.ID)
		if err != nil {
			return dto.CreateInvitationResponse{}, err
		}
		if !exists {
			toInvite = append(toInvite, u)
		}
	}

	if len(toInvite) == 0 {
		return dto.CreateInvitationResponse{}, dto.ErrInvitationAlreadyExists
	}

	// create invitation with filtered users
	inv, err := s.invitationRepo.Create(ctx, nil, entity.Invitation{EventID: eventID, Users: toInvite})
	if err != nil {
		return dto.CreateInvitationResponse{}, err
	}
	// assemble response and send emails to invited users
	names := make([]string, len(inv.Users))
	for i, u := range inv.Users { // u is entity.User
		names[i] = u.Name // This was already there

		// Fetch the UserInvitation to get the QRCode
		userInvitation, errUI := s.invitationRepo.GetUserInvitation(ctx, nil, inv.ID, u.ID)
		var pngData []byte // Declare pngData to be accessible for the final send
		var errQR error

		if errUI != nil {
			log.Printf("Error fetching UserInvitation for user %s (email: %s): %v. Sending plain invitation.", u.ID, u.Email, errUI)
			// Fallback to old SendMail or a simplified text email if QR fetch fails
			plainBody := "You have been invited to " + inv.Event.Name + ". Please contact support if you did not receive your QR code (fetch error)."
			errSend := utils.SendMail(u.Email, "Event Invitation: "+inv.Event.Name, plainBody)
			if errSend != nil {
				log.Println("Failed to send plain invitation email to", u.Email, ":", errSend)
			}
			continue // Move to next user
		}

		if userInvitation.QRCode == "" {
			log.Printf("QRCode is empty for UserInvitation UserID: %s, InvitationID: %s. Sending plain invitation.", userInvitation.UserID, userInvitation.InvitationID)
			// Fallback to old SendMail or a simplified text email if QR is empty
			plainBody := "You have been invited to " + inv.Event.Name + ". Please contact support if you did not receive your QR code (empty QR)."
			errSend := utils.SendMail(u.Email, "Event Invitation: "+inv.Event.Name, plainBody)
			if errSend != nil {
				log.Println("Failed to send plain invitation email to", u.Email, ":", errSend)
			}
			continue // Move to next user
		}

		log.Printf("Generating QR code image for UserID: %s, InvitationID: %s, QRCode: %s", u.ID, inv.ID, userInvitation.QRCode)
		pngData, errQR = qrcode.Encode(userInvitation.QRCode, qrcode.Medium, 256) // pngData is []byte
		if errQR != nil {
			log.Printf("Error generating QR code image for user %s (email: %s): %v. Sending plain invitation.", u.ID, u.Email, errQR)
			// Fallback to old SendMail or a simplified text email if QR generation fails
			plainBody := "You have been invited to " + inv.Event.Name + ". Please contact support if you did not receive your QR code (generation error)."
			errSend := utils.SendMail(u.Email, "Event Invitation: "+inv.Event.Name, plainBody)
			if errSend != nil {
				log.Println("Failed to send plain invitation email to", u.Email, ":", errSend)
			}
			continue // Move to next user
		}

		log.Printf("Successfully generated QR code image for user %s, size: %d bytes", u.ID, len(pngData))

		// Load EmailConfig to get ApiBaseUrl
		emailCfg, errCfg := config.NewEmailConfig() // This loads .env and unmarshals
		apiBaseURL := ""
		if errCfg != nil {
			log.Printf("Warning: Could not load email config to get API_BASE_URL: %v. RSVP links may be relative/broken.", errCfg)
		} else {
			apiBaseURL = emailCfg.ApiBaseUrl
		}

		if apiBaseURL == "" {
			log.Println("Warning: API_BASE_URL is not set in config. RSVP links in email will be relative or may not work as expected.")
		}

		acceptLink := apiBaseURL + "/api/invitation/rsvp/accept/" + userInvitation.QRCode
		declineLink := apiBaseURL + "/api/invitation/rsvp/decline/" + userInvitation.QRCode
		// The userInvitation.QRCode is the token

		templateData := map[string]interface{}{
			"UserName":    u.Name,
			"EventName":   inv.Event.Name,
			"Year":        time.Now().Year(),
			"AcceptLink":  acceptLink,
			"DeclineLink": declineLink,
		}

		emailSubject := "You're Invited to " + inv.Event.Name + "!"
		errSend := utils.SendInvitationMail(u.Email, emailSubject, templateData, pngData)
		if errSend != nil {
			log.Println("Failed to send styled invitation email with QR to", u.Email, ":", errSend)
			// Optional: Fallback to plain text email if styled email fails?
			// For now, just log the error and assume if SendInvitationMail fails, it's a more significant issue.
		} else {
			log.Println("Successfully sent invitation with QR code to", u.Email)
		}
	}

	now := time.Now().Format(time.RFC3339)

	return dto.CreateInvitationResponse{
		EventName:  inv.Event.Name,
		Names:      names,
		InvitedAt:  now,
		RSVPStatus: entity.RSVPStatusPending,
	}, nil
}

func (s *invitationService) GetInvitationByID(ctx context.Context, invitationID string) ([]dto.InvitationResponse, error) {
	// parse invitation ID
	id, err := uuid.Parse(invitationID)
	if err != nil {
		return nil, err
	}
	// fetch with preloads
	inv, err := s.invitationRepo.GetInvitationByID(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	// assemble list response
	resp := make([]dto.InvitationResponse, len(inv.Users))
	now := time.Now().Format(time.RFC3339)
	for i, u := range inv.Users {
		userInv, err := s.invitationRepo.GetUserInvitation(ctx, nil, inv.ID, u.ID)
		if err != nil {
			return nil, err
		}
		resp[i] = dto.InvitationResponse{
			ID:         inv.ID.String(),
			EventName:  inv.Event.Name,
			Name:       u.Name,
			InvitedAt:  now,
			RSVPStatus: entity.RSVPStatusPending,
			QRCode:     userInv.QRCode,
		}
	}
	return resp, nil
}

func (s *invitationService) GetInvitationByEventID(ctx context.Context, eventID string) ([]dto.InvitationResponse, error) {
	// parse event ID
	id, err := uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}
	// fetch invitations by event ID
	invitations, err := s.invitationRepo.GetInvitationByEvent(ctx, nil, id)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

func (s *invitationService) GetAllInvitations(ctx context.Context) ([]dto.InvitationResponse, error) {
	// TBA
	return []dto.InvitationResponse{}, dto.ErrGetAllInvitations
}

func (s *invitationService) Update(ctx context.Context, invitationID string, req dto.UpdateInvitationRequest) (dto.InvitationResponse, error) {
	// not implemented
	return dto.InvitationResponse{}, nil
}

func (s *invitationService) Delete(ctx context.Context, invitationID string) error {
	id := uuid.MustParse(invitationID)
	if err := s.invitationRepo.Delete(ctx, nil, id); err != nil {
		return err
	}
	return nil
}

// ScanQRCode handles the logic for marking attendance via QR code
func (s *invitationService) ScanQRCode(ctx context.Context, qrCode string) (dto.ScanQRCodeResponse, error) {
	userInvitation, err := s.invitationRepo.GetUserInvitationByQRCode(ctx, nil, qrCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ScanQRCodeResponse{}, errors.New("QR code not found") // Or a more specific error type/dto
		}
		return dto.ScanQRCodeResponse{}, err // Other database error
	}

	if userInvitation.AttendedAt != nil {
		return dto.ScanQRCodeResponse{}, errors.New("QR code already used") // Or a more specific error type/dto
	}

	now := time.Now()
	userInvitation.AttendedAt = &now

	updatedUserInvitation, err := s.invitationRepo.UpdateUserInvitation(ctx, nil, userInvitation)
	if err != nil {
		return dto.ScanQRCodeResponse{}, err // Error updating record
	}

	// Fetch the main Invitation record to get EventID
	inv, err := s.invitationRepo.GetInvitationByID(ctx, nil, updatedUserInvitation.InvitationID)
	if err != nil {
		// Log error, but proceed with minimal response if this secondary fetch fails
		log.Printf("Error fetching invitation details for ScanQRCode response: %v", err)
		return dto.ScanQRCodeResponse{
			UserID:     updatedUserInvitation.UserID.String(),
			AttendedAt: updatedUserInvitation.AttendedAt.Format(time.RFC3339),
			Message:    "Attendance marked successfully. Could not fetch full details.",
		}, nil
	}

	// Find the specific user in the invitation's users list
	var userName string
	for _, u := range inv.Users {
		if u.ID == updatedUserInvitation.UserID {
			userName = u.Name
			break
		}
	}
	if userName == "" {
		log.Printf("Error: User with ID %s not found in invitation %s for ScanQRCode response", updatedUserInvitation.UserID, inv.ID)
		// Attempt to get user directly if not found in the preloaded list (should not happen with correct preloading)
		// For now, we'll leave userName blank if not found in the loaded data.
	}

	return dto.ScanQRCodeResponse{
		UserID:     updatedUserInvitation.UserID.String(),
		UserName:   userName,       // This might be empty if user details are not preloaded in GetInvitationByID's Users
		EventName:  inv.Event.Name, // Relies on Event being preloaded in GetInvitationByID
		AttendedAt: updatedUserInvitation.AttendedAt.Format(time.RFC3339),
		Message:    "Attendance marked successfully",
	}, nil
}

// ProcessRSVP handles updating the RSVP status for an invitation based on a token.
func (s *invitationService) ProcessRSVP(ctx context.Context, qrCodeToken string, newRsvpStatus string) error {
	userInvitation, err := s.invitationRepo.GetUserInvitationByQRCode(ctx, nil, qrCodeToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("sorry, this RSVP link appears to be invalid or has expired")
		}
		log.Printf("Error fetching UserInvitation by QRCode token %s: %v", qrCodeToken, err)
		return errors.New("an unexpected error occurred while processing your RSVP. Please try again later")
	}

	// Check if already RSVP'd
	if userInvitation.RSVPStatus != entity.RSVPStatusPending {
		// You could customize the message further if needed, e.g. "You have already accepted this invitation on [date]."
		return errors.New("your RSVP has already been recorded as: " + userInvitation.RSVPStatus)
	}

	// Validate newRsvpStatus (though controller should send correct ones)
	if newRsvpStatus != entity.RSVPStatusAccepted && newRsvpStatus != entity.RSVPStatusDeclined {
		log.Printf("Invalid newRsvpStatus '%s' provided for token %s", newRsvpStatus, qrCodeToken)
		return errors.New("an internal error occurred. Invalid RSVP status provided") // Should not happen if called from our controller
	}

	now := time.Now()
	userInvitation.RSVPStatus = newRsvpStatus
	userInvitation.RsvpAt = &now

	_, err = s.invitationRepo.UpdateUserInvitation(ctx, nil, userInvitation)
	if err != nil {
		log.Printf("Error updating UserInvitation for token %s during RSVP: %v", qrCodeToken, err)
		return errors.New("an unexpected error occurred while saving your RSVP. Please try again later")
	}

	log.Printf("RSVP successful for token %s, new status: %s", qrCodeToken, newRsvpStatus)
	return nil // Success
}

// GetInvitationByUserID retrieves all invitations for a specific user
func (s *invitationService) GetInvitationByUserID(ctx context.Context, userID string) ([]dto.InvitationResponse, error) {
	// Parse the user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Get invitations from repository
	invitations, err := s.invitationRepo.GetInvitationByUserId(ctx, nil, uid)
	if err != nil {
		return nil, dto.ErrGetInvitationByUserID
	}
	log.Println(invitations)

	if len(invitations) == 0 {
		return []dto.InvitationResponse{}, nil
	}

	return invitations, nil
}
