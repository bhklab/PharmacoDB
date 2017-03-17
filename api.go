package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.Default()
	// version api
	v1 := router.Group("v1")
	{
		v1.GET("/cell_lines", GetCLines)
	}

	router.Run(":3000")
}
