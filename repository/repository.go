package repository

import (
	"context"

	"manifest-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryManager interface {
	DoInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	GetAgentRepo() AgentRepository
	GetVesselRepo() VesselRepository
	GetManifestRepo() ManifestRepository
	GetBC11Repo() BC11Repository
	GetNPERepo() NPERepository
}

type AgentRepository interface {
	Create(ctx context.Context, agent *models.ShippingAgent) error
	FindAll(ctx context.Context) ([]models.ShippingAgent, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.ShippingAgent, error)
}

type VesselRepository interface {
	Create(ctx context.Context, vessel *models.Vessel) error
	FindAll(ctx context.Context) ([]models.Vessel, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Vessel, error)
}

type ManifestRepository interface {
	Create(ctx context.Context, manifest *models.Manifest) error
	FindAll(ctx context.Context, page, limit int, search string) ([]models.Manifest, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Manifest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	AddDetail(ctx context.Context, detail *models.ManifestDetail) error
	Count(ctx context.Context) (int64, error)
}

type BC11Repository interface {
	Create(ctx context.Context, bc11 *models.BC11) error
	FindActiveByManifestID(ctx context.Context, manifestID uuid.UUID) (*models.BC11, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.BC11, error)
	IsContainerHasActiveBC11(ctx context.Context, containerNo string) (bool, error)
	Count(ctx context.Context) (int64, error)
}

type NPERepository interface {
	Create(ctx context.Context, npe *models.NPE) error
	Count(ctx context.Context) (int64, error)
}

type DBManager struct {
	db *gorm.DB
}

func NewDBManager(db *gorm.DB) RepositoryManager {
	return &DBManager{db: db}
}

func (m *DBManager) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return m.db.WithContext(ctx)
}

func (m *DBManager) DoInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

func (m *DBManager) GetAgentRepo() AgentRepository       { return &agentRepo{m} }
func (m *DBManager) GetVesselRepo() VesselRepository     { return &vesselRepo{m} }
func (m *DBManager) GetManifestRepo() ManifestRepository { return &manifestRepo{m} }
func (m *DBManager) GetBC11Repo() BC11Repository         { return &bc11Repo{m} }
func (m *DBManager) GetNPERepo() NPERepository           { return &npeRepo{m} }

// ---- Implementations ----

type agentRepo struct{ m *DBManager }

func (r *agentRepo) Create(ctx context.Context, a *models.ShippingAgent) error {
	return r.m.getDB(ctx).Create(a).Error
}
func (r *agentRepo) FindAll(ctx context.Context) ([]models.ShippingAgent, error) {
	var res []models.ShippingAgent
	err := r.m.getDB(ctx).Find(&res).Error
	return res, err
}
func (r *agentRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.ShippingAgent, error) {
	var res models.ShippingAgent
	err := r.m.getDB(ctx).First(&res, "id = ?", id).Error
	return &res, err
}

type vesselRepo struct{ m *DBManager }

func (r *vesselRepo) Create(ctx context.Context, v *models.Vessel) error {
	return r.m.getDB(ctx).Create(v).Error
}
func (r *vesselRepo) FindAll(ctx context.Context) ([]models.Vessel, error) {
	var res []models.Vessel
	err := r.m.getDB(ctx).Find(&res).Error
	return res, err
}
func (r *vesselRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Vessel, error) {
	var res models.Vessel
	err := r.m.getDB(ctx).First(&res, "id = ?", id).Error
	return &res, err
}

type manifestRepo struct{ m *DBManager }

func (r *manifestRepo) Create(ctx context.Context, m *models.Manifest) error {
	return r.m.getDB(ctx).Create(m).Error
}
func (r *manifestRepo) FindAll(ctx context.Context, page, limit int, search string) ([]models.Manifest, int64, error) {
	db := r.m.getDB(ctx).Model(&models.Manifest{})
	if search != "" {
		db = db.Where("manifest_number ILIKE ?", "%"+search+"%")
	}
	var total int64
	db.Count(&total)

	offset := (page - 1) * limit
	var res []models.Manifest
	err := db.Limit(limit).Offset(offset).Order("created_at desc").Find(&res).Error
	return res, total, err
}
func (r *manifestRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Manifest, error) {
	var res models.Manifest
	err := r.m.getDB(ctx).Preload("Vessel").Preload("ShippingAgent").Preload("ManifestDetails").Preload("BC11").First(&res, "id = ?", id).Error
	return &res, err
}
func (r *manifestRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.m.getDB(ctx).Model(&models.Manifest{}).Where("id = ?", id).Update("status", status).Error
}
func (r *manifestRepo) AddDetail(ctx context.Context, detail *models.ManifestDetail) error {
	return r.m.getDB(ctx).Create(detail).Error
}
func (r *manifestRepo) Count(ctx context.Context) (int64, error) {
	var c int64
	err := r.m.getDB(ctx).Model(&models.Manifest{}).Count(&c).Error
	return c, err
}

type bc11Repo struct{ m *DBManager }

func (r *bc11Repo) Create(ctx context.Context, b *models.BC11) error {
	return r.m.getDB(ctx).Create(b).Error
}
func (r *bc11Repo) FindActiveByManifestID(ctx context.Context, manifestID uuid.UUID) (*models.BC11, error) {
	var res models.BC11
	err := r.m.getDB(ctx).First(&res, "manifest_id = ? AND is_active = ?", manifestID, true).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (r *bc11Repo) FindByID(ctx context.Context, id uuid.UUID) (*models.BC11, error) {
	var res models.BC11
	err := r.m.getDB(ctx).First(&res, "id = ?", id).Error
	return &res, err
}
func (r *bc11Repo) IsContainerHasActiveBC11(ctx context.Context, containerNo string) (bool, error) {
	var count int64
	err := r.m.getDB(ctx).Table("bc11").
		Joins("JOIN manifests ON manifests.id = bc11.manifest_id").
		Joins("JOIN manifest_details ON manifest_details.manifest_id = manifests.id").
		Where("manifest_details.container_no = ? AND bc11.is_active = ?", containerNo, true).
		Count(&count).Error
	return count > 0, err
}
func (r *bc11Repo) Count(ctx context.Context) (int64, error) {
	var c int64
	err := r.m.getDB(ctx).Model(&models.BC11{}).Count(&c).Error
	return c, err
}

type npeRepo struct{ m *DBManager }

func (r *npeRepo) Create(ctx context.Context, n *models.NPE) error {
	return r.m.getDB(ctx).Create(n).Error
}
func (r *npeRepo) Count(ctx context.Context) (int64, error) {
	var c int64
	err := r.m.getDB(ctx).Model(&models.NPE{}).Count(&c).Error
	return c, err
}
