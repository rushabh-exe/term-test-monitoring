package vitals

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type res struct {
	Year string `json:"year"`
}

func Base(c *gin.Context) {
	response := []res{
		{Year: "SY"},
		{Year: "TY"},
		{Year: "LY"},
	}

	c.JSON(http.StatusOK, response)
}
