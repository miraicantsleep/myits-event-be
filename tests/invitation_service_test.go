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

// TestInvitationCreation_PopulatesQRCodeViaDBTrigger is an integration test
// to verify that the database trigger populates QRCode when an invitation is created via the service.
func TestInvitationCreation_PopulatesQRCodeViaDBTrigger(t *testing.T) {
	// 1. Setup: Real database connection and service
	db := SetUpDatabaseConnection() // From tests/db_test.go
	assert.NotNil(t, db)

	// Clean up any potential old data from previous failed runs
	// More targeted cleanup is better if specific IDs are known.
	// This is a broad cleanup for safety during testing.
	// Order matters due to foreign keys.
	db.Exec("DELETE FROM user_invitation WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'testuser.qrtrigger.%@example.com')")
	db.Exec("DELETE FROM invitations WHERE event_id IN (SELECT id FROM events WHERE name LIKE 'QR Trigger Test Event%')")
	db.Exec("DELETE FROM events WHERE name LIKE 'QR Trigger Test Event%'")
	db.Exec("DELETE FROM users WHERE email LIKE 'testuser.qrtrigger.%@example.com'")


	invitationRepo := repository.NewInvitationRepository(db)
	// For NewInvitationService, jwtService and db (*gorm.DB instance for service itself, not repo) might be needed.
	// If jwtService is not strictly needed for Create operation path, nil might be okay.
	// The db *gorm.DB for NewInvitationService is used if the service itself initiates transactions,
	// which doesn't seem to be the case for Create. Let's pass the main db connection.
	jwtService := service.NewJWTService() // A real JWT service instance
	invService := service.NewInvitationService(invitationRepo, jwtService, db)

	// 2. Create prerequisite User and Event entities directly in DB
	testUser := entity.User{
		Name:     "QR Trigger Test User",
		Email:    "testuser.qrtrigger." + uuid.NewString() + "@example.com", // Unique email
		Password: "password", // Will be hashed by BeforeCreate hook on User
		Role:     entity.RoleUser,
	}
	err := db.Create(&testUser).Error
	assert.NoError(t, err)

	testEvent := entity.Event{
		Name:        "QR Trigger Test Event " + uuid.NewString(),
		Description: "Test event for QR trigger",
		Date:        time.Now().Add(24 * time.Hour),
		Location:    "Test Location",
		Created_By:  testUser.ID,
		Type:        entity.EventTypeOnline,
		Status:      "upcoming",
	}
	err = db.Create(&testEvent).Error
	assert.NoError(t, err)

	// 3. Call invitationService.Create
	createReq := dto.CreateInvitationRequest{
		EventID: testEvent.ID.String(),
		UserIDs: []string{testUser.ID.String()},
	}
	_, err = invService.Create(context.Background(), createReq)
	assert.NoError(t, err)

	// 4. Fetch the UserInvitation record directly from the database
	var userInvitation entity.UserInvitation
	// A more robust way to get InvitationID might be needed if multiple invitations for the same event exist.
	// However, for this test, we assume one invitation is created for the event by the service call.
	// A simpler way: find the invitation first.
	var createdInvitation entity.Invitation
	errFirstInv := db.Where("event_id = ?", testEvent.ID).First(&createdInvitation).Error
	assert.NoError(t, errFirstInv)

	err = db.Where("user_id = ? AND invitation_id = ?", testUser.ID, createdInvitation.ID).First(&userInvitation).Error
	assert.NoError(t, err)


	// 5. Assert QRCode is populated and is a valid UUID
	assert.NotEmpty(t, userInvitation.QRCode, "QRCode should be populated by the database trigger")
	_, err = uuid.Parse(userInvitation.QRCode)
	assert.NoError(t, err, "Populated QRCode should be a valid UUID")

	// 6. Cleanup (delete in reverse order of creation or rely on CASCADE CONSTRAINTS if set up)
	// More specific cleanup:
	err = db.Delete(&userInvitation).Error
	assert.NoError(t, err)
	err = db.Delete(&createdInvitation).Error // This should also delete user_invitations via GORM or cascade
	assert.NoError(t, err)
	err = db.Delete(&testEvent).Error
	assert.NoError(t, err)
	err = db.Delete(&testUser).Error
	assert.NoError(t, err)
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
