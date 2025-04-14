package main

import (
	"recruitment-portal/controllers"
	"recruitment-portal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize Google OAuth
	services.InitGoogleOAuth()

	// Session middleware
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("mysession", store))

	// Templates
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/auth/login", controllers.Login)
	r.GET("/auth/callback", controllers.Callback)
	r.GET("/logout", controllers.Logout)

	r.GET("/dashboard", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(302, "/")
			return
		}
		c.JSON(200, gin.H{"message": "Logged in!", "user": user})
	})

	r.Run(":8080")
}
