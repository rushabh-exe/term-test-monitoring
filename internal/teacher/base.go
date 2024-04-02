package teacher

import "github.com/gin-gonic/gin"

func BaseGET(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GET teacher request"})
}
