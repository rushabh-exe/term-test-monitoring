package papers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	"github.com/hanshal101/term-test-monitor/helpers/auth"
)

func GetPaperRequest(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil || teacher.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding teacher"})
		return
	}

	var paperRequests []model.PaperModel
	if err = postgres.DB.Where("teacher_name = ?", teacher.Name).Find(&paperRequests).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding paper request"})
		return
	}

	c.JSON(http.StatusOK, paperRequests)
}

func CreatePaperRequest(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil || teacher.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding teacher"})
		return
	}

	var paperRequest model.PaperModel
	if err := c.BindJSON(&paperRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in json binding"})
		return
	}
	paperRequest.Request = false
	paperRequest.Status = false

	tx := postgres.DB.Begin()
	if err = tx.Create(&paperRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating paper requests"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": "created paper request"})
}

func DeletePaperRequest(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil || teacher.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding teacher"})
		return
	}

	reqID := c.Param("reqID")

	tx := postgres.DB.Begin()
	if err = tx.Where("id = ? AND teacher_name = ?", reqID, teacher.Name).Delete(&model.PaperModel{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in deleting paper request"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": "paper deleted successfully"})
}
