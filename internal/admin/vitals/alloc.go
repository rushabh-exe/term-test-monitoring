package vitals

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func UpdateMaxTeacherAlloc(c *gin.Context) {
	var ct []model.AllocationCount
	if err := c.BindJSON(&ct); err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding"})
		return
	}

	for _, cs := range ct {
		if err := postgres.DB.Model(&model.AllocationCount{}).Where("type = ?", cs.Type).Update("count", cs.Count).Error; err != nil {
			fmt.Println("error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in updating"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": "updated"})
}

func GetMaxTeacherAlloc(c *gin.Context) {
	var ct []model.AllocationCount
	if err := postgres.DB.Find(&ct).Error; err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in fetching"})
		return
	}

	c.JSON(http.StatusOK, ct)
}
