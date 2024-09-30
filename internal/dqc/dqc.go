package dqc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func GetReviews(c *gin.Context) {
	// dqcData, exists := c.Get("dqcData")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization data found"})
	// 	return
	// }

	// // Type assert the retrieved data to model.DQCMembers
	// dqc, ok := dqcData.(model.DQCMembers)
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse authorization data"})
	// 	return
	// }

	var response []model.DQCReview
	if err := postgres.DB.Find(&response).Error; err != nil {
		fmt.Printf("error: %v\n", err)
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

type desc struct {
	Description string `json:"description"`
}

func MakeReviewRequest(c *gin.Context) {
	// dqc, err := auth.GetDQC(c)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse authorization data"})
	// 	return
	// }

	var description desc
	if err := c.BindJSON(&description); err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding"})
		return
	}

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
	reviewRequest.Approver = "DQC"
	reviewRequest.Description = description.Description

	if err := tx.Save(&reviewRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in saving review requests"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": "response saved successfully"})
}
