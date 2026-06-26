package service

import (
	"context"
	"errors"

	"manifest-api/dto"
	"manifest-api/models"
	"manifest-api/repository"

	"github.com/google/uuid"
)

type AppService interface {
	CreateShippingAgent(ctx context.Context, req dto.ShippingAgentRequest) (dto.ShippingAgentResponse, error)
	GetShippingAgents(ctx context.Context) ([]dto.ShippingAgentResponse, error)
	GetShippingAgent(ctx context.Context, id string) (dto.ShippingAgentResponse, error)

	CreateVessel(ctx context.Context, req dto.VesselRequest) (dto.VesselResponse, error)
	GetVessels(ctx context.Context) ([]dto.VesselResponse, error)
	GetVessel(ctx context.Context, id string) (dto.VesselResponse, error)

	CreateManifest(ctx context.Context, req dto.ManifestRequest) (dto.ManifestResponse, error)
	GetManifests(ctx context.Context, page, limit int, search string) ([]dto.ManifestResponse, int64, error)
	GetManifest(ctx context.Context, id string) (dto.ManifestResponse, error)
	AddManifestDetail(ctx context.Context, manifestID string, req dto.ManifestDetailRequest) (dto.ManifestDetailResponse, error)

	CreateBC11(ctx context.Context, req dto.BC11Request) (dto.BC11Response, error)
	CreateNPE(ctx context.Context, req dto.NPERequest) (dto.NPEResponse, error)

	GetSummary(ctx context.Context) (dto.SummaryResponse, error)
}

type appService struct {
	repo repository.RepositoryManager
}

func NewAppService(repo repository.RepositoryManager) AppService {
	return &appService{repo: repo}
}

func (s *appService) CreateShippingAgent(ctx context.Context, req dto.ShippingAgentRequest) (dto.ShippingAgentResponse, error) {
	a := &models.ShippingAgent{
		AgentCode: req.AgentCode,
		AgentName: req.AgentName,
	}
	if err := s.repo.GetAgentRepo().Create(ctx, a); err != nil {
		return dto.ShippingAgentResponse{}, err
	}
	return dto.ShippingAgentResponse{ID: a.ID.String(), AgentCode: a.AgentCode, AgentName: a.AgentName, CreatedAt: a.CreatedAt}, nil
}

func (s *appService) GetShippingAgents(ctx context.Context) ([]dto.ShippingAgentResponse, error) {
	res, err := s.repo.GetAgentRepo().FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ShippingAgentResponse, len(res))
	for i, v := range res {
		out[i] = dto.ShippingAgentResponse{ID: v.ID.String(), AgentCode: v.AgentCode, AgentName: v.AgentName, CreatedAt: v.CreatedAt}
	}
	return out, nil
}

func (s *appService) GetShippingAgent(ctx context.Context, id string) (dto.ShippingAgentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.ShippingAgentResponse{}, errors.New("invalid uuid")
	}
	v, err := s.repo.GetAgentRepo().FindByID(ctx, uid)
	if err != nil {
		return dto.ShippingAgentResponse{}, err
	}
	return dto.ShippingAgentResponse{ID: v.ID.String(), AgentCode: v.AgentCode, AgentName: v.AgentName, CreatedAt: v.CreatedAt}, nil
}

func (s *appService) CreateVessel(ctx context.Context, req dto.VesselRequest) (dto.VesselResponse, error) {
	if req.DepartureDate.Before(req.ArrivalDate) {
		return dto.VesselResponse{}, errors.New("Departure date cannot be earlier than arrival date")
	}
	v := &models.Vessel{
		VesselCode:    req.VesselCode,
		VesselName:    req.VesselName,
		ArrivalDate:   req.ArrivalDate,
		DepartureDate: req.DepartureDate,
	}
	if err := s.repo.GetVesselRepo().Create(ctx, v); err != nil {
		return dto.VesselResponse{}, err
	}
	return dto.VesselResponse{ID: v.ID.String(), VesselCode: v.VesselCode, VesselName: v.VesselName, ArrivalDate: v.ArrivalDate, DepartureDate: v.DepartureDate, CreatedAt: v.CreatedAt}, nil
}

func (s *appService) GetVessels(ctx context.Context) ([]dto.VesselResponse, error) {
	res, err := s.repo.GetVesselRepo().FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.VesselResponse, len(res))
	for i, v := range res {
		out[i] = dto.VesselResponse{ID: v.ID.String(), VesselCode: v.VesselCode, VesselName: v.VesselName, ArrivalDate: v.ArrivalDate, DepartureDate: v.DepartureDate, CreatedAt: v.CreatedAt}
	}
	return out, nil
}

func (s *appService) GetVessel(ctx context.Context, id string) (dto.VesselResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.VesselResponse{}, errors.New("invalid uuid")
	}
	v, err := s.repo.GetVesselRepo().FindByID(ctx, uid)
	if err != nil {
		return dto.VesselResponse{}, err
	}
	return dto.VesselResponse{ID: v.ID.String(), VesselCode: v.VesselCode, VesselName: v.VesselName, ArrivalDate: v.ArrivalDate, DepartureDate: v.DepartureDate, CreatedAt: v.CreatedAt}, nil
}

func (s *appService) CreateManifest(ctx context.Context, req dto.ManifestRequest) (dto.ManifestResponse, error) {
	vid, err := uuid.Parse(req.VesselID)
	if err != nil {
		return dto.ManifestResponse{}, errors.New("invalid vessel_id")
	}
	aid, err := uuid.Parse(req.ShippingAgentID)
	if err != nil {
		return dto.ManifestResponse{}, errors.New("invalid shipping_agent_id")
	}

	m := &models.Manifest{
		ManifestNumber:  req.ManifestNumber,
		VesselID:        vid,
		ShippingAgentID: aid,
		Status:          "DRAFT",
	}
	if err := s.repo.GetManifestRepo().Create(ctx, m); err != nil {
		return dto.ManifestResponse{}, err
	}
	return dto.ManifestResponse{ID: m.ID.String(), ManifestNumber: m.ManifestNumber, VesselID: m.VesselID.String(), ShippingAgentID: m.ShippingAgentID.String(), Status: m.Status, CreatedAt: m.CreatedAt}, nil
}

func (s *appService) GetManifests(ctx context.Context, page, limit int, search string) ([]dto.ManifestResponse, int64, error) {
	res, total, err := s.repo.GetManifestRepo().FindAll(ctx, page, limit, search)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.ManifestResponse, len(res))
	for i, v := range res {
		out[i] = dto.ManifestResponse{ID: v.ID.String(), ManifestNumber: v.ManifestNumber, VesselID: v.VesselID.String(), ShippingAgentID: v.ShippingAgentID.String(), Status: v.Status, CreatedAt: v.CreatedAt}
	}
	return out, total, nil
}

func (s *appService) GetManifest(ctx context.Context, id string) (dto.ManifestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.ManifestResponse{}, errors.New("invalid uuid")
	}
	v, err := s.repo.GetManifestRepo().FindByID(ctx, uid)
	if err != nil {
		return dto.ManifestResponse{}, err
	}

	vResp := &dto.VesselResponse{ID: v.Vessel.ID.String(), VesselCode: v.Vessel.VesselCode, VesselName: v.Vessel.VesselName}
	aResp := &dto.ShippingAgentResponse{ID: v.ShippingAgent.ID.String(), AgentCode: v.ShippingAgent.AgentCode, AgentName: v.ShippingAgent.AgentName}
	var details []dto.ManifestDetailResponse
	for _, d := range v.ManifestDetails {
		details = append(details, dto.ManifestDetailResponse{ID: d.ID.String(), ManifestID: d.ManifestID.String(), ContainerNo: d.ContainerNo, GoodsDescription: d.GoodsDescription, Quantity: d.Quantity, CreatedAt: d.CreatedAt})
	}

	return dto.ManifestResponse{ID: v.ID.String(), ManifestNumber: v.ManifestNumber, VesselID: v.VesselID.String(), ShippingAgentID: v.ShippingAgentID.String(), Status: v.Status, CreatedAt: v.CreatedAt, Vessel: vResp, ShippingAgent: aResp, ManifestDetails: details}, nil
}

func (s *appService) AddManifestDetail(ctx context.Context, manifestID string, req dto.ManifestDetailRequest) (dto.ManifestDetailResponse, error) {
	mid, err := uuid.Parse(manifestID)
	if err != nil {
		return dto.ManifestDetailResponse{}, errors.New("invalid manifest_id")
	}
	
	manifest, err := s.repo.GetManifestRepo().FindByID(ctx, mid)
	if err != nil {
		return dto.ManifestDetailResponse{}, errors.New("manifest not found")
	}

	d := &models.ManifestDetail{
		ManifestID:       mid,
		ContainerNo:      req.ContainerNo,
		GoodsDescription: req.GoodsDescription,
		Quantity:         req.Quantity,
	}

	if err := s.repo.GetManifestRepo().AddDetail(ctx, d); err != nil {
		return dto.ManifestDetailResponse{}, err
	}
	
	if manifest.Status == "DRAFT" {
		s.repo.GetManifestRepo().UpdateStatus(ctx, mid, "COMPLETED")
	}

	return dto.ManifestDetailResponse{ID: d.ID.String(), ManifestID: d.ManifestID.String(), ContainerNo: d.ContainerNo, GoodsDescription: d.GoodsDescription, Quantity: d.Quantity, CreatedAt: d.CreatedAt}, nil
}

func (s *appService) CreateBC11(ctx context.Context, req dto.BC11Request) (dto.BC11Response, error) {
	mid, err := uuid.Parse(req.ManifestID)
	if err != nil {
		return dto.BC11Response{}, errors.New("invalid manifest_id")
	}

	manifest, err := s.repo.GetManifestRepo().FindByID(ctx, mid)
	if err != nil {
		return dto.BC11Response{}, errors.New("manifest not found")
	}

	for _, det := range manifest.ManifestDetails {
		isActive, err := s.repo.GetBC11Repo().IsContainerHasActiveBC11(ctx, det.ContainerNo)
		if err != nil {
			return dto.BC11Response{}, err
		}
		if isActive {
			return dto.BC11Response{}, errors.New("Container already has active BC11")
		}
	}

	bc11 := &models.BC11{
		ManifestID: mid,
		BC11Number: req.BC11Number,
		IsActive:   true,
	}

	if err := s.repo.GetBC11Repo().Create(ctx, bc11); err != nil {
		return dto.BC11Response{}, err
	}

	return dto.BC11Response{ID: bc11.ID.String(), ManifestID: bc11.ManifestID.String(), BC11Number: bc11.BC11Number, IsActive: bc11.IsActive, CreatedAt: bc11.CreatedAt}, nil
}

func (s *appService) CreateNPE(ctx context.Context, req dto.NPERequest) (dto.NPEResponse, error) {
	bc11ID, err := uuid.Parse(req.BC11ID)
	if err != nil {
		return dto.NPEResponse{}, errors.New("invalid bc11_id")
	}

	var npeResp dto.NPEResponse

	err = s.repo.DoInTransaction(ctx, func(txCtx context.Context) error {
		bc11, err := s.repo.GetBC11Repo().FindByID(txCtx, bc11ID)
		if err != nil || !bc11.IsActive {
			return errors.New("Invalid BC11")
		}

		manifest, err := s.repo.GetManifestRepo().FindByID(txCtx, bc11.ManifestID)
		if err != nil {
			return errors.New("manifest not found")
		}

		if manifest.VesselID == uuid.Nil || manifest.ShippingAgentID == uuid.Nil || len(manifest.ManifestDetails) == 0 {
			return errors.New("Manifest is not complete")
		}

		npe := &models.NPE{
			BC11ID:    bc11ID,
			NPENumber: req.NPENumber,
		}

		if err := s.repo.GetNPERepo().Create(txCtx, npe); err != nil {
			return err
		}

		if err := s.repo.GetManifestRepo().UpdateStatus(txCtx, manifest.ID, "READY_XRAY"); err != nil {
			return err
		}

		npeResp = dto.NPEResponse{ID: npe.ID.String(), BC11ID: npe.BC11ID.String(), NPENumber: npe.NPENumber, CreatedAt: npe.CreatedAt}
		return nil
	})

	if err != nil {
		return dto.NPEResponse{}, err
	}

	return npeResp, nil
}

func (s *appService) GetSummary(ctx context.Context) (dto.SummaryResponse, error) {
	manifests, _ := s.repo.GetManifestRepo().Count(ctx)
	bc11s, _ := s.repo.GetBC11Repo().Count(ctx)
	npes, _ := s.repo.GetNPERepo().Count(ctx)

	return dto.SummaryResponse{TotalManifests: manifests, TotalBC11: bc11s, TotalNPE: npes, TotalReadyXray: npes}, nil
}
