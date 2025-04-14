package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"recruitment-portal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GoogleUser struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func Login(c *gin.Context) {
	state := "random" // In production, generate secure CSRF-safe state
	url := services.GetGoogleOAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(c *gin.Context) {
	session := sessions.Default(c)

	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Code not found")
		return
	}

	token, err := services.ExchangeCodeForToken(code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange token: %s", err.Error())
		return
	}

	client := services.GoogleOAuthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get user info: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var user GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse user info")
		return
	}

	// Save to session
	session.Set("user", user)
	session.Save()

	c.Redirect(http.StatusFound, "/dashboard") // You can route to role-based logic here later
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
