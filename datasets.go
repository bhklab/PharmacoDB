package main

import "github.com/gin-gonic/gin"

// Dataset is a dataset datatype.
type Dataset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IndexDataset ...
func IndexDataset(c *gin.Context) {

}
