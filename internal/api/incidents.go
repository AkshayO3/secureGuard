package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"secureGuard/internal/data"
	"secureGuard/internal/models"
)

type IncidentHandler struct {
	IncidentModel *data.IncidentModel
}

type IncidentStruct struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	ReportedBy  string `json:"reportedBy"`
	AssignedTo  string `json:"assignedTo"`
}

type AssocIA struct {
	IID string
	AID string
}

func (i *IncidentHandler) ListIncidents(c *gin.Context) {
	incidents, err := i.IncidentModel.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

func (i *IncidentHandler) CreateIncident(c *gin.Context) {
	var incident IncidentStruct
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	inchident := &models.Incident{
		Title:       incident.Title,
		Description: incident.Description,
		Category:    incident.Category,
		Severity:    incident.Severity,
		Status:      incident.Status,
		ReportedBy:  incident.ReportedBy,
		AssignedTo:  incident.AssignedTo,
	}
	err := i.IncidentModel.Insert(inchident)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"incident": inchident})
}

func (i *IncidentHandler) GetById(c *gin.Context) {
	id := c.Param("incidentId")
	incident, err := i.IncidentModel.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"incident": incident})
}

func (i *IncidentHandler) UpdateById(c *gin.Context) {
	id := c.Param("incidentId")
	var incident IncidentStruct
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inchident := &models.Incident{
		Id:          id,
		Title:       incident.Title,
		Description: incident.Description,
		Category:    incident.Category,
		Severity:    incident.Severity,
		Status:      incident.Status,
		ReportedBy:  incident.ReportedBy,
		AssignedTo:  incident.AssignedTo,
	}
	err := i.IncidentModel.Update(inchident)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"incident": inchident})
}

func (i *IncidentHandler) DeleteById(c *gin.Context) {
	id := c.Param("incidentId")
	err := i.IncidentModel.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "incident deletion successful"})
}

func (i *IncidentHandler) AssociateIncidentWithAsset(c *gin.Context) {
	var assoc AssocIA
	if err := c.ShouldBindJSON(&assoc); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	err := i.IncidentModel.AssociateAsset(assoc.IID, assoc.AID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "incident successfully associated with asset"})
}

func (i *IncidentHandler) ListAssociatedAssets(c *gin.Context) {
	id := c.Param("incidentId")
	assets, err := i.IncidentModel.GetAssetsByIncident(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(assets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No associated assets found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"assets": assets})
}
