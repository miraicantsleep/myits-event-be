package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/utils"
	"gorm.io/gorm"
)

type (
	InvitationService interface {
		Create(ctx context.Context, req dto.CreateInvitationRequest) (dto.CreateInvitationResponse, error)
		GetInvitationByID(ctx context.Context, invitationID string) ([]dto.InvitationResponse, error)
		GetInvitationByEventID(ctx context.Context, eventID string) ([]dto.InvitationResponse, error)
		GetAllInvitations(ctx context.Context) ([]dto.InvitationResponse, error)
		Update(ctx context.Context, invitationID string, req dto.UpdateInvitationRequest) (dto.InvitationResponse, error)
		Delete(ctx context.Context, invitationID string) error
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

	// assemble response
	names := make([]string, len(inv.Users))
	for i, u := range inv.Users {
		names[i] = u.Name
	}

	now := time.Now().Format(time.RFC3339)

	// send emails
	err = utils.SendMail("mnabil190405@gmail.com", "Invitation Created", "You have been invited to an event.")
	if err != nil {
		return dto.CreateInvitationResponse{}, err
	}

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
		resp[i] = dto.InvitationResponse{
			ID:         inv.ID.String(),
			EventName:  inv.Event.Name,
			Name:       u.Name,
			InvitedAt:  now,
			RSVPStatus: entity.RSVPStatusPending,
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
