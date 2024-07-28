package vitals

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func GetTeachers(c *gin.Context) {
	teacherType := c.Param("type")
	var data []model.Main_Teachers

	tx := postgres.DB.Begin()
	switch teacherType {
	case "teachingStaff":
		// var mainTeachers []model.Main_Teachers
		if err := tx.Find(&data).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to fetch teaching staff"})
			return
		}
		// for _, teacher := range mainTeachers {
		// 	data = append(data, model.Teachers{Name: teacher.Name, Email: teacher.Email, Phone: teacher.Phone})
		// }
	case "nonteachingStaff":
		var coTeachers []model.Co_Teachers
		if err := tx.Find(&coTeachers).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to fetch non-teaching staff"})
			return
		}
		// for _, teacher := range coTeachers {
		// 	data = append(data, model.Teachers{Name: teacher.Name, Email: teacher.Email, Phone: teacher.Phone})
		// }
	default:
		c.JSON(400, gin.H{"error": "Invalid teacher type"})
		return
	}

	if len(data) < 1 {
		fmt.Println("data not present")
	}

	tx.Commit()
	c.JSON(200, data)
}

type reqTeacher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phno"`
	Type  string `json:"type"`
}

func CreateTeacher(c *gin.Context) {
	var req reqTeacher
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in Binding"})
		return
	}
	fmt.Println(req.Name)
	fmt.Println(req.Email)
	fmt.Println(req.Phone)
	fmt.Println(req.Type)
	teacher_type := req.Type
	tx := postgres.DB.Begin()
	switch teacher_type {
	case "Teaching":
		fmt.Println("entered teaching staff")
		var data model.Main_Teachers
		data.Name = req.Name
		data.Email = req.Email
		data.Phone = req.Phone
		if err := tx.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating teacher"})
			return
		}
	case "Non Teaching":
		var data model.Co_Teachers
		data.Name = req.Name
		data.Email = req.Email
		data.Phone = req.Phone
		if err := tx.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating teacher"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"error": "teacher created successfully"})
}
func EditTeacher(c *gin.Context) {
	var req reqTeacher
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in Binding"})
		return
	}
	teacherType := c.Param("type")
	tx := postgres.DB.Begin()
	switch teacherType {
	case "teachingStaff":
		var data model.Main_Teachers
		if err := tx.First(&data, "email = ?", req.Email).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
			return
		}
		data.Name = req.Name
		data.Phone = req.Phone
		if err := tx.Save(&data).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in updating teacher"})
			return
		}
	case "nonteachingStaff":
		var data model.Co_Teachers
		if err := tx.First(&data, "email = ?", req.Email).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found"})
			return
		}
		data.Name = req.Name
		data.Phone = req.Phone
		if err := tx.Save(&data).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in updating teacher"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher type"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "teacher updated successfully"})
}

func DeleteTeacher(c *gin.Context) {
	email := c.Param("email")
	teacherType := c.Param("type")

	tx := postgres.DB.Begin()

	switch teacherType {
	case "teachingStaff":
		var data model.Main_Teachers
		if err := tx.Where("email = ?", email).Delete(&data).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found or error in deletion"})
			return
		}
	case "nonteachingStaff":
		var data model.Co_Teachers
		if err := tx.Where("email = ?", email).Delete(&data).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "teacher not found or error in deletion"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher type"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "teacher deleted successfully"})
}
