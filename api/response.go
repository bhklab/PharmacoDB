package api

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderJSON outputs response as either indented or non-indented
// depending on setting by parent function.
func RenderJSON(c *gin.Context, indent bool, obj interface{}) {
	if indent {
		c.IndentedJSON(http.StatusOK, obj)
	} else {
		c.JSON(http.StatusOK, obj)
	}
}

// RenderJSONwithMeta outputs response as either indented or non-indented along with metadata about the reponse.
// Metadata includes current page, last page, per_page count and total count of result records.
func RenderJSONwithMeta(c *gin.Context, indent bool, page int, limit int, total int, include string, obj interface{}) {
	var data interface{}
	if include == "metadata" {
		lastPage := int(math.Ceil(float64(total) / float64(limit)))
		meta := gin.H{"page": page, "per_page": limit, "last_page": lastPage, "total": total}
		data = gin.H{"metadata": meta, "data": obj}
	} else {
		data = obj
	}
	if indent {
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
