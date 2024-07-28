package dqc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func GetReviews(c *gin.Context) {
	var response []model.DQCReview
	if err := postgres.DB.Find(&response).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting reviews"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetReviewbyID(c *gin.Context) {
	reqID := c.Param("reqID")
	var response model.DQCReview
	if err := postgres.DB.Where("id = ?", reqID).Find(&response).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting reviews"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func MakeReviewRequest(c *gin.Context) {
	reqID := c.Param("reqID")
	req := c.Param("req")
	var rq = false
	switch req {
	case "true":
		rq = true
	case "false":
		rq = false
	default:
		rq = false
	}

	var reviewRequest model.DQCReview
	tx := postgres.DB.Begin()
	if err := tx.Where("id = ? AND status = ? AND request = ?", reqID, false, false).Find(&reviewRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting review requests"})
		return
	}
	reviewRequest.Request = rq
	reviewRequest.Status = true

	if err := tx.Save(&reviewRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in saving review requests"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusBadRequest, gin.H{"success": "response saved successfully"})
}
