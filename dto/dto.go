package dto

import "time"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=Administrator Petugas"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type ShippingAgentRequest struct {
	AgentCode string `json:"agent_code" binding:"required"`
	AgentName string `json:"agent_name" binding:"required"`
}

type ShippingAgentResponse struct {
	ID        string    `json:"id"`
	AgentCode string    `json:"agent_code"`
	AgentName string    `json:"agent_name"`
	CreatedAt time.Time `json:"created_at"`
}

type VesselRequest struct {
	VesselCode    string    `json:"vessel_code" binding:"required"`
	VesselName    string    `json:"vessel_name" binding:"required"`
	ArrivalDate   time.Time `json:"arrival_date" binding:"required"`
	DepartureDate time.Time `json:"departure_date" binding:"required"`
}

type VesselResponse struct {
	ID            string    `json:"id"`
	VesselCode    string    `json:"vessel_code"`
	VesselName    string    `json:"vessel_name"`
	ArrivalDate   time.Time `json:"arrival_date"`
	DepartureDate time.Time `json:"departure_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type ManifestRequest struct {
	ManifestNumber  string `json:"manifest_number" binding:"required"`
	VesselID        string `json:"vessel_id" binding:"required,uuid"`
	ShippingAgentID string `json:"shipping_agent_id" binding:"required,uuid"`
}

type ManifestDetailRequest struct {
	ContainerNo      string `json:"container_no" binding:"required"`
	GoodsDescription string `json:"goods_description" binding:"required"`
	Quantity         int    `json:"quantity" binding:"required,gt=0"`
}

type ManifestDetailResponse struct {
	ID               string    `json:"id"`
	ManifestID       string    `json:"manifest_id"`
	ContainerNo      string    `json:"container_no"`
	GoodsDescription string    `json:"goods_description"`
	Quantity         int       `json:"quantity"`
	CreatedAt        time.Time `json:"created_at"`
}

type ManifestResponse struct {
	ID              string                   `json:"id"`
	ManifestNumber  string                   `json:"manifest_number"`
	VesselID        string                   `json:"vessel_id"`
	ShippingAgentID string                   `json:"shipping_agent_id"`
	Status          string                   `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	Vessel          *VesselResponse          `json:"vessel,omitempty"`
	ShippingAgent   *ShippingAgentResponse   `json:"shipping_agent,omitempty"`
	ManifestDetails []ManifestDetailResponse `json:"manifest_details,omitempty"`
}

type BC11Request struct {
	ManifestID string `json:"manifest_id" binding:"required,uuid"`
	BC11Number string `json:"bc11_number" binding:"required"`
}

type BC11Response struct {
	ID         string    `json:"id"`
	ManifestID string    `json:"manifest_id"`
	BC11Number string    `json:"bc11_number"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type NPERequest struct {
	BC11ID    string `json:"bc11_id" binding:"required,uuid"`
	NPENumber string `json:"npe_number" binding:"required"`
}

type NPEResponse struct {
	ID        string    `json:"id"`
	BC11ID    string    `json:"bc11_id"`
	NPENumber string    `json:"npe_number"`
	CreatedAt time.Time `json:"created_at"`
}

type SummaryResponse struct {
	TotalManifests int64 `json:"total_manifests"`
	TotalBC11      int64 `json:"total_bc11"`
	TotalNPE       int64 `json:"total_npe"`
	TotalReadyXray int64 `json:"total_ready_xray"`
}
