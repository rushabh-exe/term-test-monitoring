package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
)

func TeacherAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("teacherData")
		fmt.Println("Cookie : ", cookie)
		if err != nil {
			// fmt.Printf("error : %v", err)
			fmt.Println("cookie not present")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cookie fetch error (middleware)"})
			return
		}
		decodedData, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			fmt.Println("Error decoding base64 data:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Decoding error (middleware)"})
			return
		}
		var teacher model.Main_Teachers
		if err := json.Unmarshal([]byte(decodedData), &teacher); err != nil {
			fmt.Fprintf(os.Stderr, "Error : %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid cookie data (middleware)"})
			return
		}
		c.Set("teacherData", teacher)

		c.Next()
	}
}

func DQCAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("dqcData")
		if err != nil {
			fmt.Printf("error : %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cookie fetch error"})
			return
		}

		decodedData, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			fmt.Println("Error decoding base64 data:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Decoding error"})
			return
		}

		var dqc model.DQCMembers
		if err := json.Unmarshal([]byte(decodedData), &dqc); err != nil {
			fmt.Fprintf(os.Stderr, "Error : %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid cookie data"})
			return
		}

		c.Set("dqcData", dqc)
		c.Next()
	}
}
