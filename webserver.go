package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func launchWebserver(debugMode bool) {

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.LoadHTMLGlob("templates/*") // load html templates

	router.GET("/", indexPage)
	router.POST("/", indexPagePOST)
	router.StaticFS("/static", http.Dir("static"))

	api_routes := router.Group("/api")
	{
		//api_routes.GET("/count", apiGetCategoryCount)
		api_routes.GET("/list/:category", apiGetCategoryList)
		//api_routes.GET("/info/:itemID", apiGetItemInfo)
		//api_routes.POST("/import", apiImportNewItem)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "pong",
			"timestamp": time.Now().Unix(), //for status checking
		})
	})
	router.Run(":80") // http://localhost:1337
}

func indexPage(c *gin.Context) {

	/*//TODO: error handling/clean this up
	encodedHash, _ := c.Request.Cookie("access")
	decodedHash, _ := url.QueryUnescape(encodedHash.Value)
	isAuthed, _ := userIsAuthenticated("default", []byte(decodedHash))

	if !isAuthed {
		fmt.Println("Page hit w/o password.")
		c.HTML(http.StatusOK, "templates/index_locked.html", gin.H{
			"title": "Please authorize :)",
		})
	} else { // user authed ok
		c.HTML(http.StatusOK, "templates/index_unlocked.html", gin.H{})
	}*/

	c.HTML(http.StatusOK, "templates/index_locked.html", gin.H{
		"title": "Login",
	})
}

// check the password provided on the front page.
func indexPagePOST(c *gin.Context) {

	/*
		hashedInput, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("key")), bcrypt.DefaultCost)

		authed, hash := userIsAuthenticated("default", hashedInput)

		if authed { // user entered the correct password.
			fmt.Println("tried to set cookie")
			c.SetCookie("access", url.QueryEscape(string(hash)), 3600, "/", "localhost", true, true)
		}
	*/

	if c.PostForm("key") == "password123" {
		c.HTML(http.StatusOK, "templates/index_unlocked.html", gin.H{})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/") // reload page
	}

}
