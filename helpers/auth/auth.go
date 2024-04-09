package auth

// import (
// 	"encoding/base64"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hanshal101/term-test-monitor/database/model"
// 	"github.com/hanshal101/term-test-monitor/database/postgres"
// )

// type User struct {
// 	Name     string
// 	Email    string
// 	Password string
// 	Role     string
// }

// func IsAdmin(c *gin.Context) {
// 	userRole := getUserContext()
// 	if userRole.Role != "ADMIN" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Access"})
// 		c.Abort()
// 		return
// 	}
// }

// func getUserContext() User {
// 	user := User{
// 		Name:     "Hanshal",
// 		Email:    "lol@mail.com",
// 		Password: "123",
// 		Role:     "TEACHER",
// 	}
// 	return user
// }

// type Req struct {
// 	Email string `json:"email"`
// }

// type Res struct {
// 	Status int
// 	Error  string
// }

// func IsTeacher(c *gin.Context) {
// 	var req Req
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in binding request"})
// 	}
// 	emailExists := false
// 	var data []model.Main_Teachers
// 	if err := postgres.DB.Find(&data).Error; err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error in Main_Teacher models"})
// 		return
// 	}
// 	var teacher model.Main_Teachers
// 	for _, teachers := range data {
// 		if teachers.Email == req.Email {
// 			emailExists = true
// 			teacher.Email = teachers.Email
// 			teacher.Name = teachers.Name
// 			teacher.Phone = teachers.Phone
// 			return
// 		}
// 	}
// 	encode := base64.StdEncoding.EncodeToString()
// 	print(encode)
// 	if !emailExists {
// 		res := Res{
// 			Status: 401,
// 			Error:  "User UnAuthorized",
// 		}
// 		c.JSON(http.StatusUnauthorized, res)
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK)
// }
