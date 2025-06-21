package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"secureGuard/internal/data"
	"secureGuard/internal/models"
)

type AssetHandler struct {
	AssetModel *data.AssetModel
}

type AssetObject struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	IpAddr string `json:"ipAddr"`
	Os     string `json:"os"`
	Status string `json:"status"`
}

func (a *AssetHandler) List(c *gin.Context) {
	assets, err := a.AssetModel.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"assets": assets})
}

func (a *AssetHandler) Get(c *gin.Context) {
	id := c.Param("id")
	asset, err := a.AssetModel.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"asset": asset})
}

func (a *AssetHandler) Insert(c *gin.Context) {
	var ast AssetObject
	if err := c.ShouldBindJSON(&ast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	asset := &models.Asset{
		Name:   ast.Name,
		Type:   ast.Type,
		Ip:     ast.IpAddr,
		Os:     ast.Os,
		Status: ast.Status,
	}
	if err := a.AssetModel.Insert(asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"asset": asset})
}

func (a *AssetHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var asset AssetObject
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetX := &models.Asset{
		Id:     id,
		Name:   asset.Name,
		Type:   asset.Type,
		Ip:     asset.IpAddr,
		Os:     asset.Os,
		Status: asset.Status,
	}
	log.Println(assetX)
	err := a.AssetModel.Update(assetX)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"asset": assetX})
}

func (a *AssetHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := a.AssetModel.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset deletion successful"})
}

func (a *AssetHandler) GetAssociatedV(c *gin.Context) {
	id := c.Param("id")
	vulns, err := a.AssetModel.GetVulnerabilities(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(vulns) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No associated vulnerabilities found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"associatedVulnerabilities": vulns})
}

func (a *AssetHandler) GetAssociatedI(c *gin.Context) {
	id := c.Param("id")
	incidents, err := a.AssetModel.GetIncidents(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(incidents) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No associated incidents found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"associatedIncidents": incidents})
}
