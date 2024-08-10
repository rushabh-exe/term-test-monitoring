package dqc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	"github.com/hanshal101/term-test-monitor/helpers/auth"
)

func GetReviewRequest(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil || teacher.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding teacher"})
		return
	}

	var response []model.DQCReview
	if err = postgres.DB.Where("name = ? AND email = ?", teacher.Name, teacher.Email).Find(&response).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting reviews"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func CreateDQCReview(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil || teacher.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding teacher"})
		return
	}

	var req model.DQCReview
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding json"})
		return
	}

	tx := postgres.DB.Begin()
	req.Name = teacher.Name
	req.Email = teacher.Email
	req.Request = false
	req.Status = false
	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating review"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": "successfully snet the review"})
}

func DeleteDQCReview(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil || teacher.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding teacher"})
		return
	}

	reqID := c.Param("reqID")

	tx := postgres.DB.Begin()
	if err := tx.Where("id = ? AND name = ? AND email = ?", reqID, teacher.Name, teacher.Email).Delete(&model.DQCReview{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in deletign review"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": "successfully deleted the review"})
}
