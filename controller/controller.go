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

func (ctrl *AppController) CreateShippingAgent(c *gin.Context) {
	var req dto.ShippingAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateShippingAgent(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

func (ctrl *AppController) GetShippingAgents(c *gin.Context) {
	res, err := ctrl.svc.GetShippingAgents(c.Request.Context())
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) GetShippingAgent(c *gin.Context) {
	res, err := ctrl.svc.GetShippingAgent(c.Request.Context(), c.Param("id"))
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) CreateVessel(c *gin.Context) {
	var req dto.VesselRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateVessel(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

func (ctrl *AppController) GetVessels(c *gin.Context) {
	res, err := ctrl.svc.GetVessels(c.Request.Context())
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) GetVessel(c *gin.Context) {
	res, err := ctrl.svc.GetVessel(c.Request.Context(), c.Param("id"))
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) CreateManifest(c *gin.Context) {
	var req dto.ManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateManifest(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

func (ctrl *AppController) GetManifests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	res, _, err := ctrl.svc.GetManifests(c.Request.Context(), page, limit, search)
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) GetManifest(c *gin.Context) {
	res, err := ctrl.svc.GetManifest(c.Request.Context(), c.Param("id"))
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) AddManifestDetail(c *gin.Context) {
	var req dto.ManifestDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.AddManifestDetail(c.Request.Context(), c.Param("id"), req)
	handleResponse(c, res, err, http.StatusCreated)
}

func (ctrl *AppController) CreateBC11(c *gin.Context) {
	var req dto.BC11Request
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateBC11(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

func (ctrl *AppController) CreateNPE(c *gin.Context) {
	var req dto.NPERequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, nil, err, http.StatusBadRequest)
		return
	}
	res, err := ctrl.svc.CreateNPE(c.Request.Context(), req)
	handleResponse(c, res, err, http.StatusCreated)
}

func (ctrl *AppController) GetSummary(c *gin.Context) {
	res, err := ctrl.svc.GetSummary(c.Request.Context())
	handleResponse(c, res, err, http.StatusOK)
}

func (ctrl *AppController) SeedData(c *gin.Context) {
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Seed triggered successfully",
	})
}
