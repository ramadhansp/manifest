package controller

import (
	"net/http"
	"strconv"

	"manifest-api/dto"
	"manifest-api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppController struct {
	svc service.AppService
}

func NewAppController(svc service.AppService) *AppController {
	return &AppController{svc: svc}
}

func handleResponse(c *gin.Context, data interface{}, err error, status int) {
	if err != nil {
		var errs []string
		if valErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range valErrs {
				errs = append(errs, e.Field()+" is "+e.Tag())
			}
		} else {
			errs = append(errs, err.Error())
		}
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  errs,
		})
		return
	}
	c.JSON(status, dto.APIResponse{
		Success: true,
		Message: "Success",
		Data:    data,
	})
}

// Register godoc
// @Summary     Registrasi user baru
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body dto.RegisterRequest true "Data registrasi"
// @Success     201 {object} dto.APIResponse{data=dto.RegisterResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /register [post]
func (ctrl *AppController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.Register(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// Login godoc
// @Summary     Login user
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body dto.LoginRequest true "Kredensial login"
// @Success     200 {object} dto.APIResponse{data=dto.LoginResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /login [post]
func (ctrl *AppController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.Login(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusOK)
}

// CreateShippingAgent godoc
// @Summary     Tambah shipping agent
// @Tags        Shipping Agent
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.ShippingAgentRequest true "Data shipping agent"
// @Success     201 {object} dto.APIResponse{data=dto.ShippingAgentResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /shipping-agents [post]
func (ctrl *AppController) CreateShippingAgent(c *gin.Context) {
	var req dto.ShippingAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateShippingAgent(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// GetShippingAgents godoc
// @Summary     List semua shipping agent
// @Tags        Shipping Agent
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} dto.APIResponse{data=[]dto.ShippingAgentResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /shipping-agents [get]
func (ctrl *AppController) GetShippingAgents(c *gin.Context) {
	res, err := ctrl.svc.GetShippingAgents(c.Request.Context())
	handleResponse(c, res, err, http.StatusOK)
}

// GetShippingAgent godoc
// @Summary     Detail shipping agent by ID
// @Tags        Shipping Agent
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "Shipping Agent ID (UUID)"
// @Success     200 {object} dto.APIResponse{data=dto.ShippingAgentResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /shipping-agents/{id} [get]
func (ctrl *AppController) GetShippingAgent(c *gin.Context) {
	res, err := ctrl.svc.GetShippingAgent(c.Request.Context(), c.Param("id"))
	handleResponse(c, res, err, http.StatusOK)
}

// CreateVessel godoc
// @Summary     Tambah vessel (kapal)
// @Tags        Vessel
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.VesselRequest true "Data vessel"
// @Success     201 {object} dto.APIResponse{data=dto.VesselResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /vessels [post]
func (ctrl *AppController) CreateVessel(c *gin.Context) {
	var req dto.VesselRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateVessel(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// GetVessels godoc
// @Summary     List semua vessel
// @Tags        Vessel
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} dto.APIResponse{data=[]dto.VesselResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /vessels [get]
func (ctrl *AppController) GetVessels(c *gin.Context) {
	res, err := ctrl.svc.GetVessels(c.Request.Context())
	handleResponse(c, res, err, http.StatusOK)
}

// GetVessel godoc
// @Summary     Detail vessel by ID
// @Tags        Vessel
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "Vessel ID (UUID)"
// @Success     200 {object} dto.APIResponse{data=dto.VesselResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /vessels/{id} [get]
func (ctrl *AppController) GetVessel(c *gin.Context) {
	res, err := ctrl.svc.GetVessel(c.Request.Context(), c.Param("id"))
	handleResponse(c, res, err, http.StatusOK)
}

// CreateManifest godoc
// @Summary     Buat manifest baru
// @Tags        Manifest
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.ManifestRequest true "Data manifest"
// @Success     201 {object} dto.APIResponse{data=dto.ManifestResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /manifests [post]
func (ctrl *AppController) CreateManifest(c *gin.Context) {
	var req dto.ManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateManifest(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// GetManifests godoc
// @Summary     List manifest dengan pagination dan filter
// @Tags        Manifest
// @Security    BearerAuth
// @Produce     json
// @Param       page   query int    false "Halaman (default: 1)"
// @Param       limit  query int    false "Jumlah per halaman (default: 10)"
// @Param       search query string false "Cari berdasarkan nomor manifest"
// @Param       id     query string false "Filter berdasarkan ID manifest (UUID)"
// @Success     200 {object} dto.APIResponse{data=[]dto.ManifestResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /manifests [get]
func (ctrl *AppController) GetManifests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	id := c.Query("id")
	res, _, err := ctrl.svc.GetManifests(c.Request.Context(), page, limit, search, id)
	handleResponse(c, res, err, http.StatusOK)
}

// GetManifest godoc
// @Summary     Detail manifest by ID (include vessel, agent, containers, BC11)
// @Tags        Manifest
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "Manifest ID (UUID)"
// @Success     200 {object} dto.APIResponse{data=dto.ManifestResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /manifests/{id} [get]
func (ctrl *AppController) GetManifest(c *gin.Context) {
	res, err := ctrl.svc.GetManifest(c.Request.Context(), c.Param("id"))
	handleResponse(c, res, err, http.StatusOK)
}

// AddManifestDetail godoc
// @Summary     Tambah detail container ke manifest
// @Tags        Manifest
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id      path string                   true "Manifest ID (UUID)"
// @Param       request body dto.ManifestDetailRequest true "Data container"
// @Success     201 {object} dto.APIResponse{data=dto.ManifestDetailResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /manifests/{id}/details [post]
func (ctrl *AppController) AddManifestDetail(c *gin.Context) {
	var req dto.ManifestDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.AddManifestDetail(c.Request.Context(), c.Param("id"), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// CreateBC11 godoc
// @Summary     Buat BC11 untuk manifest
// @Tags        BC11
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.BC11Request true "Data BC11"
// @Success     201 {object} dto.APIResponse{data=dto.BC11Response}
// @Failure     400 {object} dto.APIResponse
// @Router      /bc11 [post]
func (ctrl *AppController) CreateBC11(c *gin.Context) {
	var req dto.BC11Request
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateBC11(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// CreateNPE godoc
// @Summary     Buat NPE berdasarkan BC11
// @Tags        NPE
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.NPERequest true "Data NPE"
// @Success     201 {object} dto.APIResponse{data=dto.NPEResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /npe [post]
func (ctrl *AppController) CreateNPE(c *gin.Context) {
	var req dto.NPERequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateNPE(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

// GetSummary godoc
// @Summary     Ringkasan statistik data
// @Tags        Summary
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} dto.APIResponse{data=dto.SummaryResponse}
// @Failure     400 {object} dto.APIResponse
// @Router      /summary [get]
func (ctrl *AppController) GetSummary(c *gin.Context) {
	res, err := ctrl.svc.GetSummary(c.Request.Context())
	handleResponse(c, res, err, http.StatusOK)
}
