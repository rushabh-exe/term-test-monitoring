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

func GetDQC(c *gin.Context) (model.DQCMembers, error) {
	var dqc model.DQCMembers

	cookie, err := c.Cookie("dqcData")
	if err != nil {
		fmt.Printf("error : %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cookie fetch error"})
		return dqc, err
	}

	fmt.Println("Cookie content:", cookie)

	decodedData, err := base64.StdEncoding.DecodeString(cookie)
	if err != nil {
		fmt.Println("Error decoding base64 data:", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Decoding error"})
		return dqc, err
	}

	if err := json.Unmarshal([]byte(decodedData), &dqc); err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid cookie data"})
		return dqc, err
	}

	return dqc, nil
}
