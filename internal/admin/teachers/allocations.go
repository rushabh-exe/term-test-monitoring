package teachers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func CreateTeacherAllocation(c *gin.Context) {
	fetch, err := http.Get("http://localhost:3002/api/allotment")
	if err != nil {
		log.Fatalf("Error in getting Teacher Allocation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher allocation"})
		return
	}
	defer fetch.Body.Close()

	var result []model.TeacherAllocation
	if err := json.NewDecoder(fetch.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body"})
		return
	}

	tx := postgres.DB.Begin()
	if err := tx.Where("deleted_at IS NOT NULL").Delete(&model.TeacherAllocation{}).Error; err != nil {
		log.Fatalf("Error in DB Deletion: %v", err)
	}

	for _, item := range result {
		teacherAllocation := model.TeacherAllocation{
			Date:         item.Date,
			Start_Time:   item.Start_Time,
			End_Time:     item.End_Time,
			Classroom:    item.Classroom,
			Main_Teacher: item.Main_Teacher,
			Co_Teacher:   item.Co_Teacher,
		}
		if err := tx.Create(&teacherAllocation).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Error in DB Creation: %v", err)
		}
	}
	tx.Commit()

	// Process result further as needed...

	c.JSON(http.StatusOK, gin.H{"status": "Allocation Done"})
}

func GetTeacherAllocation(c *gin.Context) {
	var teacherAllocations []model.TeacherAllocation
	if err := postgres.DB.Where("deleted_at IS NULL").Find(&teacherAllocations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher allocations"})
		return
	}

	c.JSON(http.StatusOK, teacherAllocations)
}

func DeleteTeacherAllocation(c *gin.Context) {
	id := c.Param("id")

	var data model.TeacherAllocation

	tx := postgres.DB.Begin()
	if err := tx.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in deleting DB"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"error": "Allocation Deleted successfully"})
}
