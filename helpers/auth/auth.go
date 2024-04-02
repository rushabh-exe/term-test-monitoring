package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func IsAdmin(c *gin.Context) {
	userRole := getUserContext()
	if userRole.Role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Access"})
		c.Abort()
		return
	}
}

func getUserContext() User {
	user := User{
		Name:     "Hanshal",
		Email:    "lol@mail.com",
		Password: "123",
		Role:     "TEACHER",
	}
	return user
}

type Req struct {
	Name  string
	Email string
}

type Res struct {
	Status int
	Error  string
}

func IsTeacher(c *gin.Context) {
	var req Req
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in binding request"})
	}
	var emails = []string{
		"lol@mail.com",
		"s.mishra@mail.com",
		"j.khanapuri@mail.com",
		"test@mail.com",
	}

	emailExists := false
	for _, email := range emails {
		if email == req.Email {
			emailExists = true
			break
		}
	}
	if emailExists {
		res := Res{
			Status: 200,
			Error:  "User Authenticated",
		}
		c.JSON(http.StatusOK, res)
		return
	}
	res := Res{
		Status: 401,
		Error:  "User UnAuthorized",
	}
	c.JSON(http.StatusUnauthorized, res)
	c.Abort()
}
