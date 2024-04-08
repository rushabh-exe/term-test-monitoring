package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// TeacherData represents the structure of the teacher data stored in the cookie
type TeacherData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func TeacherAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("teacherData")
		if err != nil {
			fmt.Printf("error : %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cookie fetch error"})
			return
		}
		decodedData, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			fmt.Println("Error decoding base64 data:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Decoding error"})
			return
		}
		var teacher TeacherData
		if err := json.Unmarshal([]byte(decodedData), &teacher); err != nil {
			fmt.Fprintf(os.Stderr, "Error : %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid cookie data"})
			return
		}
		c.Set("teacherData", teacher)

		c.Next()
	}
}
