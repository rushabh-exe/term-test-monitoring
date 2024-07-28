package vitals

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

type req struct {
	Subject string `json:"subject"`
}

func GetSubject(c *gin.Context) {
	year := c.Param("year")
	var res []model.Subject
	tx := postgres.DB.Begin()
	if err := tx.Where("year = ?", year).Find(&res).Error; err != nil {
		fmt.Fprintf(os.Stderr, "Error : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting subjects"})
	}
	tx.Commit()
	c.JSON(http.StatusOK, res)
}

func CreateSubject(c *gin.Context) {
	year := c.Param("year")
	var subject model.Subject
	var Req req
	if err := c.BindJSON(&Req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding"})
		return
	}
	if Req.Subject == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in subject"})
		return
	}
	subject = model.Subject{
		Year: year,
		Name: Req.Subject,
	}
	tx := postgres.DB.Begin()
	if err := tx.Create(&subject).Error; err != nil {
		fmt.Fprintf(os.Stderr, "Error : ", err)
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"error": "subject created"})
}

func DeleteSubject(c *gin.Context) {
	year := c.Param("year")
	subjectName := c.Param("subject")

	tx := postgres.DB.Begin()
	if err := tx.Where("year = ? AND name = ?", year, subjectName).Delete(&model.Subject{}).Error; err != nil {
		tx.Rollback()
		fmt.Fprintf(os.Stderr, "Error deleting subject: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subject"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Subject deleted successfully"})
}
