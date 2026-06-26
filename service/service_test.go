package service

import (
	"context"
	"testing"

	"manifest-api/dto"
	"manifest-api/models"
	"manifest-api/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepoManager struct {
	mock.Mock
	userRepo     *MockUserRepo
	agentRepo    *MockAgentRepo
	vesselRepo   *MockVesselRepo
	manifestRepo *MockManifestRepo
	bc11Repo     *MockBC11Repo
	npeRepo      *MockNPERepo
}

func (m *MockRepoManager) DoInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
func (m *MockRepoManager) GetUserRepo() repository.UserRepository         { return m.userRepo }
func (m *MockRepoManager) GetAgentRepo() repository.AgentRepository       { return m.agentRepo }
func (m *MockRepoManager) GetVesselRepo() repository.VesselRepository     { return m.vesselRepo }
func (m *MockRepoManager) GetManifestRepo() repository.ManifestRepository { return m.manifestRepo }
func (m *MockRepoManager) GetBC11Repo() repository.BC11Repository         { return m.bc11Repo }
func (m *MockRepoManager) GetNPERepo() repository.NPERepository           { return m.npeRepo }

type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) Create(ctx context.Context, user *models.User) error {
	return m.Called(ctx, user).Error(0)
}

type MockAgentRepo struct{ mock.Mock }

func (m *MockAgentRepo) Create(ctx context.Context, a *models.ShippingAgent) error {
	return m.Called(ctx, a).Error(0)
}
func (m *MockAgentRepo) FindAll(ctx context.Context) ([]models.ShippingAgent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.ShippingAgent), args.Error(1)
}
func (m *MockAgentRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.ShippingAgent, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ShippingAgent), args.Error(1)
}

type MockVesselRepo struct{ mock.Mock }

func (m *MockVesselRepo) Create(ctx context.Context, v *models.Vessel) error {
	return m.Called(ctx, v).Error(0)
}
func (m *MockVesselRepo) FindAll(ctx context.Context) ([]models.Vessel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Vessel), args.Error(1)
}
func (m *MockVesselRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Vessel, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Vessel), args.Error(1)
}

type MockManifestRepo struct{ mock.Mock }

func (m *MockManifestRepo) Create(ctx context.Context, mn *models.Manifest) error {
	return m.Called(ctx, mn).Error(0)
}
func (m *MockManifestRepo) FindAll(ctx context.Context, page, limit int, search string) ([]models.Manifest, int64, error) {
	args := m.Called(ctx, page, limit, search)
	return args.Get(0).([]models.Manifest), args.Get(1).(int64), args.Error(2)
}
func (m *MockManifestRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Manifest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Manifest), args.Error(1)
}
func (m *MockManifestRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockManifestRepo) AddDetail(ctx context.Context, d *models.ManifestDetail) error {
	return m.Called(ctx, d).Error(0)
}
func (m *MockManifestRepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockBC11Repo struct{ mock.Mock }

func (m *MockBC11Repo) Create(ctx context.Context, b *models.BC11) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockBC11Repo) FindActiveByManifestID(ctx context.Context, id uuid.UUID) (*models.BC11, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.BC11), args.Error(1)
}
func (m *MockBC11Repo) FindByID(ctx context.Context, id uuid.UUID) (*models.BC11, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BC11), args.Error(1)
}
func (m *MockBC11Repo) IsContainerHasActiveBC11(ctx context.Context, containerNo string) (bool, error) {
	args := m.Called(ctx, containerNo)
	return args.Get(0).(bool), args.Error(1)
}
func (m *MockBC11Repo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockNPERepo struct{ mock.Mock }

func (m *MockNPERepo) Create(ctx context.Context, n *models.NPE) error {
	return m.Called(ctx, n).Error(0)
}
func (m *MockNPERepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateBC11_Success(t *testing.T) {
	mockManRepo := new(MockManifestRepo)
	mockBC11Repo := new(MockBC11Repo)
	mgr := &MockRepoManager{manifestRepo: mockManRepo, bc11Repo: mockBC11Repo}
	svc := NewAppService(mgr)

	manifestID := uuid.New()
	manifest := &models.Manifest{
		ID: manifestID,
		ManifestDetails: []models.ManifestDetail{
			{ContainerNo: "CONT123"},
		},
	}

	mockManRepo.On("FindByID", mock.Anything, manifestID).Return(manifest, nil)
	mockBC11Repo.On("IsContainerHasActiveBC11", mock.Anything, "CONT123").Return(false, nil)
	mockBC11Repo.On("Create", mock.Anything, mock.AnythingOfType("*models.BC11")).Return(nil)

	req := dto.BC11Request{
		ManifestID: manifestID.String(),
		BC11Number: "BC11-001",
	}

	res, err := svc.CreateBC11(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "BC11-001", res.BC11Number)
	assert.True(t, res.IsActive)
	mockManRepo.AssertExpectations(t)
	mockBC11Repo.AssertExpectations(t)
}

func TestCreateBC11_FailContainerHasActiveBC11(t *testing.T) {
	mockManRepo := new(MockManifestRepo)
	mockBC11Repo := new(MockBC11Repo)
	mgr := &MockRepoManager{manifestRepo: mockManRepo, bc11Repo: mockBC11Repo}
	svc := NewAppService(mgr)

	manifestID := uuid.New()
	manifest := &models.Manifest{
		ID: manifestID,
		ManifestDetails: []models.ManifestDetail{
			{ContainerNo: "CONT123"},
		},
	}

	mockManRepo.On("FindByID", mock.Anything, manifestID).Return(manifest, nil)
	// Return true here means container already has active BC11
	mockBC11Repo.On("IsContainerHasActiveBC11", mock.Anything, "CONT123").Return(true, nil)

	req := dto.BC11Request{
		ManifestID: manifestID.String(),
		BC11Number: "BC11-001",
	}

	_, err := svc.CreateBC11(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "Container already has active BC11", err.Error())
	mockManRepo.AssertExpectations(t)
	mockBC11Repo.AssertExpectations(t)
}

func TestCreateNPE_Success(t *testing.T) {
	mockManRepo := new(MockManifestRepo)
	mockBC11Repo := new(MockBC11Repo)
	mockNPERepo := new(MockNPERepo)
	mgr := &MockRepoManager{manifestRepo: mockManRepo, bc11Repo: mockBC11Repo, npeRepo: mockNPERepo}
	svc := NewAppService(mgr)

	bc11ID := uuid.New()
	manifestID := uuid.New()

	bc11 := &models.BC11{
		ID:         bc11ID,
		ManifestID: manifestID,
		IsActive:   true,
	}

	manifest := &models.Manifest{
		ID:              manifestID,
		VesselID:        uuid.New(), // Complete manifest
		ShippingAgentID: uuid.New(), // Complete manifest
		ManifestDetails: []models.ManifestDetail{
			{ContainerNo: "CONT123"},
		},
	}

	mockBC11Repo.On("FindByID", mock.Anything, bc11ID).Return(bc11, nil)
	mockManRepo.On("FindByID", mock.Anything, manifestID).Return(manifest, nil)
	mockNPERepo.On("Create", mock.Anything, mock.AnythingOfType("*models.NPE")).Return(nil)
	mockManRepo.On("UpdateStatus", mock.Anything, manifestID, "READY_XRAY").Return(nil)

	req := dto.NPERequest{
		BC11ID:    bc11ID.String(),
		NPENumber: "NPE-001",
	}

	res, err := svc.CreateNPE(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "NPE-001", res.NPENumber)
	mockBC11Repo.AssertExpectations(t)
	mockManRepo.AssertExpectations(t)
	mockNPERepo.AssertExpectations(t)
}

func TestCreateNPE_FailManifestIncomplete(t *testing.T) {
	mockManRepo := new(MockManifestRepo)
	mockBC11Repo := new(MockBC11Repo)
	mgr := &MockRepoManager{manifestRepo: mockManRepo, bc11Repo: mockBC11Repo}
	svc := NewAppService(mgr)

	bc11ID := uuid.New()
	manifestID := uuid.New()

	bc11 := &models.BC11{
		ID:         bc11ID,
		ManifestID: manifestID,
		IsActive:   true,
	}

	manifest := &models.Manifest{
		ID:              manifestID,
		VesselID:        uuid.Nil, // Incomplete manifest
		ShippingAgentID: uuid.Nil,
		ManifestDetails: []models.ManifestDetail{},
	}

	mockBC11Repo.On("FindByID", mock.Anything, bc11ID).Return(bc11, nil)
	mockManRepo.On("FindByID", mock.Anything, manifestID).Return(manifest, nil)

	req := dto.NPERequest{
		BC11ID:    bc11ID.String(),
		NPENumber: "NPE-001",
	}

	_, err := svc.CreateNPE(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "Manifest is not complete", err.Error())
	mockBC11Repo.AssertExpectations(t)
	mockManRepo.AssertExpectations(t)
}

func TestCreateNPE_FailInvalidBC11(t *testing.T) {
	mockBC11Repo := new(MockBC11Repo)
	mgr := &MockRepoManager{bc11Repo: mockBC11Repo}
	svc := NewAppService(mgr)

	bc11ID := uuid.New()

	// BC11 found but inactive
	bc11 := &models.BC11{
		ID:       bc11ID,
		IsActive: false,
	}

	mockBC11Repo.On("FindByID", mock.Anything, bc11ID).Return(bc11, nil)

	req := dto.NPERequest{
		BC11ID:    bc11ID.String(),
		NPENumber: "NPE-001",
	}

	_, err := svc.CreateNPE(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "Invalid BC11", err.Error())
	mockBC11Repo.AssertExpectations(t)
}
