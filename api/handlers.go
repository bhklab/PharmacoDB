package api

import "github.com/gin-gonic/gin"

// IndexCells ...
func IndexCells(c *gin.Context) {
	var cells Cells
	err := cells.List()
	if err != nil {
		c.String(404, "nothing found")
	}
	RenderJSON(c, true, cells)
}
