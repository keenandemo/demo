package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Media struct {
	ID       int    `db:"id",json:"id"`
	Title    string `db:"title",json:"title"`
	ThumbSrc string `db:"thumbnailPath",json:"thumbSrc"`
	IMDB_ID  string `db:"imdb_id",json:"imdb_id"`
}

type Category struct {
	List []Media `json:"list"`
}

func apiGetCategoryList(c *gin.Context) {
	var category string

	switch c.Param("category") {
	case "movies":
		category = "movies"
	case "shows":
		category = "shows"
	default:
		category = "movies"
	}

	listOfMedia := new(Category)

	mediaTable := dbcon.Collection(category)
	result := mediaTable.Find()

	err := result.All(&listOfMedia.List)

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status": "error",
			"data":   "{}",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"data":   listOfMedia,
	})

}
