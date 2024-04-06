package alloc_helper

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func GetTeachers(c *gin.Context) {
	var ttalloc []model.CreateTimeTable
	var mainTeachers []model.Main_Teachers
	var coTeachers []model.Co_Teachers

	if err := postgres.DB.Find(&mainTeachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch main teachers"})
		return
	}

	if err := postgres.DB.Find(&coTeachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch co teachers"})
		return
	}

	if err := postgres.DB.Find(&ttalloc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tt alloc"})
		return
	}

	response := map[string]interface{}{
		"main_teachers": mainTeachers,
		"co_teachers":   coTeachers,
		"ttalloc":       ttalloc,
	}

	c.JSON(http.StatusOK, response)
}

func GetClass(c *gin.Context) {
	var response []model.AllocationResult
	if err := postgres.DB.Select("DISTINCT class_room").Find(&response).Error; err != nil {
		log.Fatalf("Error in extraction of alloc: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocation results"})
		return
	}

	classrooms := make([]string, len(response))
	for i, res := range response {
		classrooms[i] = res.ClassRoom
	}

	c.JSON(http.StatusOK, gin.H{"classroom": classrooms})

}
