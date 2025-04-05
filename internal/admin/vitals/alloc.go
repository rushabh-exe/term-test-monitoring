package vitals

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

type AllocationCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type AllocationPayload struct {
	Count string `json:"count"`
	Type  string `json:"type"`
}

func UpdateMaxTeacherAlloc(c *gin.Context) {
	var payload []AllocationPayload
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println("error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	for _, item := range payload {
		countInt, err := strconv.Atoi(item.Count)
		if err != nil {
			fmt.Printf("Invalid count value '%s' for type '%s': %v\n", item.Count, item.Type, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid count '%s' for type '%s'", item.Count, item.Type)})
			return
		}
		condition := model.AllocationCount{Type: item.Type}
		dataToAssignOrCreate := model.AllocationCount{Count: strconv.Itoa(countInt)}

		result := postgres.DB.Where(condition).
			Assign(dataToAssignOrCreate).
			FirstOrCreate(&model.AllocationCount{})

		if result.Error != nil {
			fmt.Printf("Error upserting allocation count for type %s: %v\n", item.Type, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error processing type " + item.Type})
			return
		}
		fmt.Printf("Processed type '%s', rows affected: %d\n", item.Type, result.RowsAffected)
	}

	c.JSON(http.StatusOK, gin.H{"success": "Allocation counts updated successfully"})
}

func GetMaxTeacherAlloc(c *gin.Context) {
	var ct []AllocationCount
	var ct_db []model.AllocationCount
	if err := postgres.DB.Order("type asc").Find(&ct_db).Error; err != nil {
		fmt.Println("error fetching allocation counts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching allocation counts"})
		return
	}
	for _, item := range ct_db {
		countInt, err := strconv.Atoi(item.Count)
		if err != nil {
			fmt.Println("error in converting string to int")
		}
		ct = append(ct, AllocationCount{Type: item.Type, Count: countInt})
	}

	if ct == nil {
		ct = []AllocationCount{}
	}

	c.JSON(http.StatusOK, ct)
}
