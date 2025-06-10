package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockInvitationRepository is a mock implementation of InvitationRepository
type MockInvitationRepository struct {
	mock.Mock
}

func (m *MockInvitationRepository) Create(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error) {
	args := m.Called(ctx, tx, invitation)
	return args.Get(0).(entity.Invitation), args.Error(1)
}

func (m *MockInvitationRepository) GetInvitationByID(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) (entity.Invitation, error) {
	args := m.Called(ctx, tx, invitationID)
	// Ensure that the first argument is of type entity.Invitation or can be asserted to it.
	// If it can be nil or another type on error, handle that.
	var inv entity.Invitation
	if args.Get(0) != nil {
		inv = args.Get(0).(entity.Invitation)
	}
	return inv, args.Error(1)
}

func (m *MockInvitationRepository) GetInvitationByEvent(ctx context.Context, tx *gorm.DB, eventID uuid.UUID) ([]dto.InvitationResponse, error) {
	args := m.Called(ctx, tx, eventID)
	// Similar handling for nil or type assertion
	var resp []dto.InvitationResponse
	if args.Get(0) != nil {
		resp = args.Get(0).([]dto.InvitationResponse)
	}
	return resp, args.Error(1)
}

func (m *MockInvitationRepository) GetAllUserInvitations(ctx context.Context, tx *gorm.DB) ([]entity.Invitation, error) {
	args := m.Called(ctx, tx)
	var invs []entity.Invitation
	if args.Get(0) != nil {
		invs = args.Get(0).([]entity.Invitation)
	}
	return invs, args.Error(1)
}

func (m *MockInvitationRepository) Update(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error) {
	args := m.Called(ctx, tx, invitation)
	var inv entity.Invitation
	if args.Get(0) != nil {
		inv = args.Get(0).(entity.Invitation)
	}
	return inv, args.Error(1)
}

func (m *MockInvitationRepository) Delete(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) error {
	args := m.Called(ctx, tx, invitationID)
	return args.Error(0)
}

func (m *MockInvitationRepository) CheckInvitationExist(ctx context.Context, tx *gorm.DB, eventID uuid.UUID, userID ...uuid.UUID) (bool, error) {
	// Pass userID as a separate argument to m.Called for varargs
	allArgs := []interface{}{ctx, tx, eventID}
	for _, id := range userID {
		allArgs = append(allArgs, id)
	}
	args := m.Called(allArgs...)
	return args.Bool(0), args.Error(1)
}

func (m *MockInvitationRepository) CheckInvitationRSVPStatus(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID, rsvpStatus string) (bool, error) {
	args := m.Called(ctx, tx, invitationID, rsvpStatus)
	return args.Bool(0), args.Error(1)
}

func (m *MockInvitationRepository) GetUserInvitationByQRCode(ctx context.Context, tx *gorm.DB, qrCode string) (entity.UserInvitation, error) {
	args := m.Called(ctx, tx, qrCode)
	var userInv entity.UserInvitation
	if args.Get(0) != nil {
		userInv = args.Get(0).(entity.UserInvitation)
	}
	return userInv, args.Error(1)
}

func (m *MockInvitationRepository) UpdateUserInvitation(ctx context.Context, tx *gorm.DB, userInvitation entity.UserInvitation) (entity.UserInvitation, error) {
	args := m.Called(ctx, tx, userInvitation)
	var userInv entity.UserInvitation
	if args.Get(0) != nil {
		userInv = args.Get(0).(entity.UserInvitation)
	}
	return userInv, args.Error(1)
}

func TestInvitationService_ScanQRCode_Success(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	// No need for jwtService or db for this specific service method test if fully mocked
	invService := service.NewInvitationService(mockRepo, nil, nil)

	qrCode := "valid-qr-code"
	testUserID := uuid.New()
	testInvitationID := uuid.New()
	testEventID := uuid.New()
	now := time.Now()

	userInv := entity.UserInvitation{
		UserID:       testUserID,
		InvitationID: testInvitationID,
		QRCode:       qrCode,
		AttendedAt:   nil, // Not attended yet
	}

	// Expected data for GetInvitationByID
	expectedEvent := entity.Event{ID: testEventID, Name: "Test Event"}
	expectedUser := entity.User{ID: testUserID, Name: "Test User"}
	expectedInvitation := entity.Invitation{
		ID:      testInvitationID,
		EventID: testEventID,
		Event:   expectedEvent,
		Users:   []entity.User{expectedUser},
	}

	mockRepo.On("GetUserInvitationByQRCode", mock.Anything, (*gorm.DB)(nil), qrCode).Return(userInv, nil)
	// We need to capture the argument to UpdateUserInvitation to check AttendedAt
	mockRepo.On("UpdateUserInvitation", mock.Anything, (*gorm.DB)(nil), mock.AnythingOfType("entity.UserInvitation")).Run(func(args mock.Arguments) {
		updatedUserInv := args.Get(2).(entity.UserInvitation)
		assert.NotNil(t, updatedUserInv.AttendedAt)
	}).Return(func(ctx context.Context, db *gorm.DB, ui entity.UserInvitation) entity.UserInvitation {
		ui.AttendedAt = &now // Simulate setting the time
		return ui
	}, nil)
	mockRepo.On("GetInvitationByID", mock.Anything, (*gorm.DB)(nil), testInvitationID).Return(expectedInvitation, nil)


	response, err := invService.ScanQRCode(context.Background(), qrCode)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, testUserID.String(), response.UserID)
	assert.Equal(t, "Test User", response.UserName)
	assert.Equal(t, "Test Event", response.EventName)
	assert.NotEmpty(t, response.AttendedAt)
	assert.Equal(t, "Attendance marked successfully", response.Message)
	mockRepo.AssertExpectations(t)
}

func TestInvitationService_ScanQRCode_NotFound(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	invService := service.NewInvitationService(mockRepo, nil, nil)
	qrCode := "invalid-qr-code"

	mockRepo.On("GetUserInvitationByQRCode", mock.Anything, (*gorm.DB)(nil), qrCode).Return(entity.UserInvitation{}, gorm.ErrRecordNotFound)

	response, err := invService.ScanQRCode(context.Background(), qrCode)

	assert.Error(t, err)
	assert.Equal(t, "QR code not found", err.Error())
	assert.Empty(t, response)
	mockRepo.AssertExpectations(t)
}

func TestInvitationService_ScanQRCode_AlreadyUsed(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	invService := service.NewInvitationService(mockRepo, nil, nil)
	qrCode := "used-qr-code"
	attendedTime := time.Now().Add(-1 * time.Hour)

	userInv := entity.UserInvitation{
		QRCode:     qrCode,
		AttendedAt: &attendedTime, // Already attended
	}

	mockRepo.On("GetUserInvitationByQRCode", mock.Anything, (*gorm.DB)(nil), qrCode).Return(userInv, nil)

	response, err := invService.ScanQRCode(context.Background(), qrCode)

	assert.Error(t, err)
	assert.Equal(t, "QR code already used", err.Error())
	assert.Empty(t, response)
	mockRepo.AssertExpectations(t)
}

func TestInvitationService_ScanQRCode_UpdateError(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	invService := service.NewInvitationService(mockRepo, nil, nil)
	qrCode := "valid-qr-code-for-update-fail"
    testUserID := uuid.New()
	testInvitationID := uuid.New()

	userInv := entity.UserInvitation{
        UserID:       testUserID,
		InvitationID: testInvitationID,
		QRCode:       qrCode,
		AttendedAt:   nil,
	}
	dbError := errors.New("database update error")

	mockRepo.On("GetUserInvitationByQRCode", mock.Anything, (*gorm.DB)(nil), qrCode).Return(userInv, nil)
	mockRepo.On("UpdateUserInvitation", mock.Anything, (*gorm.DB)(nil), mock.AnythingOfType("entity.UserInvitation")).Return(entity.UserInvitation{}, dbError)

	response, err := invService.ScanQRCode(context.Background(), qrCode)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	assert.Empty(t, response)
	mockRepo.AssertExpectations(t)
}

func TestInvitationService_ScanQRCode_GetInvitationDetailsError(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	invService := service.NewInvitationService(mockRepo, nil, nil)

	qrCode := "valid-qr-code-details-fail"
	testUserID := uuid.New()
	testInvitationID := uuid.New()
	now := time.Now()

	userInv := entity.UserInvitation{
		UserID:       testUserID,
		InvitationID: testInvitationID,
		QRCode:       qrCode,
		AttendedAt:   nil,
	}

	// Mock successful GetUserInvitationByQRCode and UpdateUserInvitation
	mockRepo.On("GetUserInvitationByQRCode", mock.Anything, (*gorm.DB)(nil), qrCode).Return(userInv, nil)
	mockRepo.On("UpdateUserInvitation", mock.Anything, (*gorm.DB)(nil), mock.AnythingOfType("entity.UserInvitation")).Return(func(ctx context.Context, db *gorm.DB, ui entity.UserInvitation) entity.UserInvitation {
		ui.AttendedAt = &now
		return ui
	}, nil)

	// Mock GetInvitationByID to return an error
	mockRepo.On("GetInvitationByID", mock.Anything, (*gorm.DB)(nil), testInvitationID).Return(entity.Invitation{}, errors.New("failed to get invitation details"))

	response, err := invService.ScanQRCode(context.Background(), qrCode)

	assert.NoError(t, err) // The main operation succeeded, error is in fetching details for response
	assert.NotNil(t, response)
	assert.Equal(t, testUserID.String(), response.UserID)
	assert.Empty(t, response.UserName) // UserName and EventName might be empty or have default values
	assert.Empty(t, response.EventName)
	assert.NotEmpty(t, response.AttendedAt)
	// The message should indicate partial success
	// The current service code logs the error but returns a specific message if inv details fetch fails.
    assert.Equal(t, "Attendance marked successfully. Could not fetch full details.", response.Message)


	mockRepo.AssertExpectations(t)
}

```
