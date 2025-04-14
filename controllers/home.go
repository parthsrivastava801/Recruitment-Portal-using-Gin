package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"title":   "Recruitment Portal",
		"message": "Welcome to the Recruitment Portal!",
	})
}
