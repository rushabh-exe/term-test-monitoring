package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

type authResp struct {
	Cookie string `json:"cookie"`
}

func IsTeacherAuth(c *gin.Context) {
	var req model.EAuthReq
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

	c.JSON(http.StatusOK, authResp{
		Cookie: encode,
	})
}

func GetTeacher(c *gin.Context) (model.Main_Teachers, error) {
	var teacher model.Main_Teachers

	cookie, err := c.Request.Cookie("teacherData")
	if err != nil {
		fmt.Printf("error : %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cookie fetch error"})
	}

	decodedData, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		fmt.Println("Error decoding base64 data:", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Decoding error"})
	}

	if err := json.Unmarshal([]byte(decodedData), &teacher); err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid cookie data"})
	}

	return teacher, nil
}
