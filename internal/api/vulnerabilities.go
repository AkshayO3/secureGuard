package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"secureGuard/internal/data"
	"secureGuard/internal/models"
)

type VulnerabilityHandler struct {
	VulnerabilityModel *data.VulnerabilityModel
}

type VulnerabilityObject struct {
	CveId       string
	Title       string
	Description string
	Severity    string
	Status      string
}

type AssocVA struct {
	VID    string
	AID    string
	Status string
}

func (v *VulnerabilityHandler) List(c *gin.Context) {
	vulns, err := v.VulnerabilityModel.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(vulns) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No vulnerabilities found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vulnerabilities": vulns})
}

func (v *VulnerabilityHandler) Insert(c *gin.Context) {
	var vuln VulnerabilityObject
	if err := c.ShouldBindJSON(&vuln); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	vulnR := &models.Vulnerability{
		CveId:       vuln.CveId,
		Title:       vuln.Title,
		Description: vuln.Description,
		Severity:    vuln.Severity,
		Status:      vuln.Status,
	}
	err := v.VulnerabilityModel.Insert(vulnR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"vulnerability": vulnR})
}

func (v *VulnerabilityHandler) Get(c *gin.Context) {
	id := c.Param("id")
	vuln, err := v.VulnerabilityModel.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vulnerability": vuln})
}

func (v *VulnerabilityHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var vuln VulnerabilityObject
	if err := c.ShouldBindJSON(&vuln); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vulnR := &models.Vulnerability{
		Id:          id,
		CveId:       vuln.CveId,
		Title:       vuln.Title,
		Description: vuln.Description,
		Severity:    vuln.Severity,
		Status:      vuln.Status,
	}
	err := v.VulnerabilityModel.Update(vulnR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vulnerability": vulnR})
}

func (v *VulnerabilityHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := v.VulnerabilityModel.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "vulnerability successfully deleted"})
}

func (v *VulnerabilityHandler) GetAssociatedAssets(c *gin.Context) {
	id := c.Param("vulnId")
	assets, err := v.VulnerabilityModel.GetAssetsByVulnerability(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(assets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No associated assets found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"associatedAssets": assets})
}

func (v *VulnerabilityHandler) AssociateA(c *gin.Context) {
	var obj AssocVA
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := v.VulnerabilityModel.AssociateAsset(obj.VID, obj.AID, obj.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset associated successfully"})
}
