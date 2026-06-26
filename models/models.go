package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShippingAgent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	AgentCode string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	AgentName string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Vessel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	VesselCode    string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	VesselName    string    `gorm:"type:varchar(100);not null"`
	ArrivalDate   time.Time `gorm:"not null"`
	DepartureDate time.Time `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Manifest struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	ManifestNumber  string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	VesselID        uuid.UUID `gorm:"type:uuid;not null"`
	ShippingAgentID uuid.UUID `gorm:"type:uuid;not null"`
	Status          string    `gorm:"type:varchar(20);not null;default:'DRAFT'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Vessel          Vessel           `gorm:"foreignKey:VesselID"`
	ShippingAgent   ShippingAgent    `gorm:"foreignKey:ShippingAgentID"`
	ManifestDetails []ManifestDetail `gorm:"foreignKey:ManifestID"`
	BC11            *BC11            `gorm:"foreignKey:ManifestID"`
}

type ManifestDetail struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	ManifestID       uuid.UUID `gorm:"type:uuid;not null"`
	ContainerNo      string    `gorm:"type:varchar(50);not null"`
	GoodsDescription string    `gorm:"type:text;not null"`
	Quantity         int       `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type BC11 struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ManifestID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	BC11Number string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	IsActive   bool      `gorm:"not null;default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	NPE *NPE `gorm:"foreignKey:BC11ID"`
}

type NPE struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	BC11ID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	NPENumber string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (base *ShippingAgent) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}
func (base *Vessel) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}
func (base *Manifest) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	if base.Status == "" {
		base.Status = "DRAFT"
	}
	return nil
}
func (base *ManifestDetail) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}
func (base *BC11) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}
func (base *NPE) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}
