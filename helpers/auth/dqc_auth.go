package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func IsDqcAuth(c *gin.Context) {
	var req model.EAuthReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in Binding"})
		return
	}
	var data model.DQCMembers
	if err := postgres.DB.First(&data, "email = ?", req.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dqc not found"})
		return
	}
	jsonData, err := json.Marshal(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in marshal"})
		return
	}
	encode := base64.StdEncoding.EncodeToString(jsonData)

	c.JSON(http.StatusOK, authResp{
		Cookie: encode,
	})
}
