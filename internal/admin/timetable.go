package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func CreateTimeTable(c *gin.Context) {
	year := c.Param("year")
	sem := c.Param(":sem")
	var req []model.CreateTimeTable

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tx := postgres.DB.Begin()
	for _, tt := range req {
		tt.Year = year
		tt.Sem = sem
		if err := tx.Create(&tt).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create timetable"})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Timetable created successfully"})
}

func GetTT(c *gin.Context) {
	var timetables []model.CreateTimeTable

	if err := postgres.DB.Find(&timetables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"timetables": timetables})
}

func GetTTbyYear(c *gin.Context) {
	year := c.Param("Year")

	var timetables []model.CreateTimeTable

	if err := postgres.DB.Where("year = ?", year).Find(&timetables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"timetables": timetables})
}

func DeleteTimeTable(c *gin.Context) {
	year := c.Param("year")

	var timetable []model.CreateTimeTable

	tx := postgres.DB.Begin()
	if err := tx.Where("year = ?", year).Delete(&timetable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Delete timetables"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"error": "Timetable deleted successfully"})
}
