package students

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	allocate "github.com/hanshal101/term-test-monitor/helpers/Allocation"
)

func DualAllocation(c *gin.Context) {
	var reqArr model.ReqArr
	if err := c.BindJSON(&reqArr); err != nil {
		log.Fatalf("Error in Binding")
	}

	result := []model.AllocationResult{}

	for _, req := range reqArr.ReqAll {
		totalCapacity := 0
		for _, cap := range req.Class {
			totalCapacity += int(cap.Capacity)
		}

		max_no := math.Max(float64(req.Class1.Strength), float64(req.Class2.Strength))
		if max_no > float64(totalCapacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Space"})
		}

		roomRollNo := 1
		tx := postgres.DB.Begin()

		for i, class := range req.Class {
			result = append(result, model.AllocationResult{
				ClassRoom: class.Room,
				ClassName: req.Class1.Name,
				Start:     int64(roomRollNo),
				End:       allocate.Allot(req.Class1.Strength, class.Capacity, i),
			})
			result = append(result, model.AllocationResult{
				ClassRoom: class.Room,
				ClassName: req.Class2.Name,
				Start:     int64(roomRollNo),
				End:       allocate.Allot(req.Class2.Strength, class.Capacity, i),
			})

			roomRollNo += int(class.Capacity)
		}

		if err := tx.Create(&result).Error; err != nil {
			tx.Rollback()
			fmt.Printf("Error : %v\n", err)
			return
		}

		tx.Commit()
	}

	c.JSON(http.StatusOK, result)
}

func GetAllocation(c *gin.Context) {
	var response []model.AllocationResult
	if err := postgres.DB.Find(&response).Error; err != nil {
		log.Fatalf("Error in extraction of alloc: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocation results"})
		return
	}

	print(response)

	c.JSON(http.StatusOK, response)
}

func SingleAllocation(c *gin.Context) {
	var reqArr model.SingleAllocReqArr
	if err := c.BindJSON(&reqArr); err != nil {
		log.Fatalf("Error in Binding")
	}

	result := []model.AllocationResult{}

	for _, req := range reqArr.ReqAll {
		totalCapacity := 0
		for _, cap := range req.Class {
			totalCapacity += int(cap.Capacity)
		}

		if float64(req.Class1.Strength) > float64(totalCapacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Space"})
		}

		roomRollNo := 1
		tx := postgres.DB.Begin()

		for i, class := range req.Class {
			result = append(result, model.AllocationResult{
				ClassRoom: class.Room,
				ClassName: req.Class1.Name,
				Start:     int64(roomRollNo),
				End:       allocate.Allot(req.Class1.Strength, class.Capacity, i),
			})

			roomRollNo += int(class.Capacity)
		}

		if err := tx.Create(&result).Error; err != nil {
			tx.Rollback()
			fmt.Printf("Error : %v\n", err)
			return
		}

		tx.Commit()
	}

	c.JSON(http.StatusOK, result)
}
