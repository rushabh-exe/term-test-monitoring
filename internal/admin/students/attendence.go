package students

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func DashboardAttendence(c *gin.Context) {
	year := c.Param("year")
	subject := c.Param("subject")
	class := c.Param("class")

	var response []model.Attendence_Models
	tx := postgres.DB.Begin()
	if err := tx.Where("year = ? AND subject = ? AND class = ?", year, subject, class).Find(&response).Error; err != nil {
		fmt.Fprintf(os.Stderr, "Error : \n", err)
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func DeleteAttendence(c *gin.Context) {
	year := c.Param("year")
	subject := c.Param("subject")
	class := c.Param("class")

	tx := postgres.DB.Begin()
	if err := tx.Where("year = ? AND subject = ? AND class = ?", year, subject, class).Delete(&model.Attendence_Models{}).Error; err != nil {
		tx.Rollback()
		fmt.Fprintf(os.Stderr, "Error deleting attendance: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendance"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Attendance deleted successfully"})
}

func EditAttendence(c *gin.Context) {
	var attendanceReq []model.Attendence_Models
	if err := c.BindJSON(&attendanceReq); err != nil {
		fmt.Fprintf(os.Stderr, "Error in binding: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "Error occurred: %v\n", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	for _, req := range attendanceReq {
		conditions := &model.Attendence_Models{
			Name:       req.Name,
			RollNo:     req.RollNo,
			Year:       req.Year,
			Class:      req.Class,
			Subject:    req.Subject,
			M_Teacher:  req.M_Teacher,
			C_Teacher:  req.C_Teacher,
			Date:       req.Date,
			Start_Time: req.Start_Time,
			End_Time:   req.End_Time,
		}

		if err := tx.Model(&model.Attendence_Models{}).Where(conditions).Update("is_present", req.IsPresent).Error; err != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "Error updating record: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Attendance records updated successfully"})
}
