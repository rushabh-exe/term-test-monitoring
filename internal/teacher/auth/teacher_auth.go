package teacher_auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

type authResp struct {
	TeacherData string `json:"teacherData"`
}

type authReq struct {
	Email string `json:"email"`
}

func Auth(c *gin.Context) {
	var req authReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in Binding"})
		return
	}
	var data model.Main_Teachers
	if err := postgres.DB.First(&data, "email = ?", req.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
		return
	}
	jsonData, err := json.Marshal(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in marshal"})
		return
	}
	encode := base64.StdEncoding.EncodeToString(jsonData)

	resp := authResp{
		TeacherData: encode,
	}
	c.JSON(http.StatusOK, resp)
}
