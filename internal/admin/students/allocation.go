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

// func DualAllocation(c *gin.Context) {
// 	var reqArr []model.DualAllocationReq
// 	if err := c.BindJSON(&reqArr); err != nil {
// 		log.Fatalf("Error in Binding")
// 	}

// 	result := []model.AllocationResult{}

// 	for _, req := range reqArr {
// 		totalCapacity := 0
// 		for _, cap := range req.Class {
// 			totalCapacity += int(cap.Capacity)
// 		}

// 		max_no := math.Max(float64(req.Class1.Strength), float64(req.Class2.Strength))
// 		if max_no > float64(totalCapacity) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Space"})
// 		}

// 		roomRollNo := 1
// 		tx := postgres.DB.Begin()

// 		for i, class := range req.Class {
// 			result = append(result, model.AllocationResult{
// 				ClassRoom: class.Room,
// 				ClassName: req.Class1.Name,
// 				Start:     int64(roomRollNo),
// 				End:       allocate.Allot(req.Class1.Strength, class.Capacity, i),
// 			})
// 			result = append(result, model.AllocationResult{
// 				ClassRoom: class.Room,
// 				ClassName: req.Class2.Name,
// 				Start:     int64(roomRollNo),
// 				End:       allocate.Allot(req.Class2.Strength, class.Capacity, i),
// 			})

// 			roomRollNo += int(class.Capacity)
// 		}

// 		if err := tx.Create(&result).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Printf("Error : %v\n", err)
// 			return
// 		}

// 		tx.Commit()
// 	}

//		c.JSON(http.StatusOK, result)
//	}

// 50, 30, 1
func Allocate(strength int, capacity int, index int) int64 {
	allocated := min(strength, capacity)
	if index*capacity+allocated > strength {
		return int64(strength)
	}
	return int64(index*capacity + allocated)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func DualAllocation(c *gin.Context) {
	var reqArr []model.DualAllocationReq
	if err := c.BindJSON(&reqArr); err != nil {
		log.Fatalf("Error in Binding: %v", err)
	}

	var results []model.AllocationResult

	for _, req := range reqArr {
		totalCapacity := 0
		for _, class := range req.Class {
			totalCapacity += int(class.Capacity)
		}

		maxStrength := math.Max(float64(req.Class1.Strength), float64(req.Class2.Strength))
		if maxStrength > float64(totalCapacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Space"})
			return
		}

		classARollNo, classBRollNo := 1, 1
		tx := postgres.DB.Begin()

		for i, class := range req.Class {
			classAEnd := Allocate(int(req.Class1.Strength), int(class.Capacity), i)
			if classARollNo <= int(req.Class1.Strength) {
				results = append(results, model.AllocationResult{
					ClassRoom: class.Room,
					ClassName: req.Class1.Name,
					Start:     int64(classARollNo),
					End:       classAEnd,
				})
				classARollNo += min(int(req.Class1.Strength), int(class.Capacity))
			} else {
				results = append(results, model.AllocationResult{
					ClassRoom: class.Room,
					ClassName: req.Class1.Name,
					Start:     0,
					End:       0,
				})
			}

			classBEnd := Allocate(int(req.Class2.Strength), int(class.Capacity), i)
			if classBRollNo <= int(req.Class2.Strength) {
				results = append(results, model.AllocationResult{
					ClassRoom: class.Room,
					ClassName: req.Class2.Name,
					Start:     int64(classBRollNo),
					End:       classBEnd,
				})
				classBRollNo += min(int(req.Class2.Strength), int(class.Capacity))
			} else {
				results = append(results, model.AllocationResult{
					ClassRoom: class.Room,
					ClassName: req.Class2.Name,
					Start:     0,
					End:       0,
				})
			}
		}

		if err := tx.Create(&results).Error; err != nil {
			tx.Rollback()
			fmt.Printf("Error : %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save allocation results"})
			return
		}

		tx.Commit()
	}

	c.JSON(http.StatusOK, results)
}

func GetAllocation(c *gin.Context) {
	var response []model.AllocationResult
	if err := postgres.DB.Find(&response).Error; err != nil {
		log.Fatalf("Error in extraction of alloc: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocation results"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func SingleAllocation(c *gin.Context) {
	var reqArr []model.SingleAllocReq
	if err := c.BindJSON(&reqArr); err != nil {
		log.Fatalf("Error in Binding")
	}

	result := []model.AllocationResult{}

	for _, req := range reqArr {
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

func DeleteAllocation(c *gin.Context) {
	id := c.Param("id")
	var data model.AllocationResult
	tx := postgres.DB.Begin()
	if err := tx.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in deleting DB"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"error": "Allocation Deleted"})
}
