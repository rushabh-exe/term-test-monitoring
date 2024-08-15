package teachers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func GetPaperRequests(c *gin.Context) {
	var response []model.PaperModel
	if err := postgres.DB.Find(&response).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting paper requests"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetPaperRequestsStatus(c *gin.Context) {
	var response []model.PaperModel
	if err := postgres.DB.Where("status = ?", false).Find(&response).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting paper requests"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func MakePaperRequests(c *gin.Context) {
	reqID := c.Param("reqID")
	req := c.Param("req")
	var rq bool
	switch req {
	case "true":
		rq = true
	case "false":
		rq = false
	default:
		rq = false
	}

	var paperRequest model.PaperModel
	tx := postgres.DB.Begin()
	if err := tx.Where("id = ? AND status = ? AND request = ?", reqID, false, false).Find(&paperRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting paper requests"})
		return
	}
	paperRequest.Request = rq
	paperRequest.Status = true

	if err := tx.Save(&paperRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in saving paper requests"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": "response saved successfully"})
}
