package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Synonym is a match between a datatype name and datasets that use the name.
type Synonym struct {
	Name     string    `json:"name"`
	Datasets []Dataset `json:"datasets"`
}

// Set Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://71d8d1bc8e4843eeba979fdaadebe48b:df30d2048fc44b5185809f04ba9d2294@sentry.io/186627")
}

// Prepare database abstraction for later use.
func initDB() (*sql.DB, error) {
	cred := os.Getenv("mysql_user") + ":" + os.Getenv("mysql_passwd") + "@tcp(127.0.0.1:3306)/pharmacodb"
	db, err := sql.Open("mysql", cred)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	// Establish and test database connection.
	err = db.Ping()
	if err != nil {
		raven.CaptureError(err, nil)
	}
	return db, err
}

// Handle error messages (all except no route match errors).
func handleError(c *gin.Context, err error, code int, message string) {
	if err != nil {
		raven.CaptureError(err, nil)
	}
	c.JSON(code, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// writeHeaderLinks writes pagination links in response header.
// Links available under 'Link' header, including (prev, next, first, last).
func writeHeaderLinks(c *gin.Context, endpoint string, page int, total int, limit int) {
	var (
		prev    string
		prevRel string
		next    string
		nextRel string
	)
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	first := fmt.Sprintf("<https://api.pharmacodb.com/v1%s?page=%d&per_page=%d>", endpoint, 1, limit)
	if (page > 1) && (page <= lastPage) {
		prev = fmt.Sprintf("<https://api.pharmacodb.com/v1%s?page=%d&per_page=%d>", endpoint, page-1, limit)
		prevRel = "; rel=\"prev\", "
	}
	if (page >= 1) && (page < lastPage) {
		next = fmt.Sprintf("<https://api.pharmacodb.com/v1%s?page=%d&per_page=%d>", endpoint, page+1, limit)
		nextRel = "; rel=\"next\", "
	}
	last := fmt.Sprintf("<https://api.pharmacodb.com/v1%s?page=%d&per_page=%d>", endpoint, lastPage, limit)

	linknp := prev + prevRel + next + nextRel
	linkfl := first + "; rel=\"first\", " + last + "; rel=\"last\""
	link := linknp + linkfl

	c.Writer.Header().Set("Link", link)
}
