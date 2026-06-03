package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zept888/go-project-278/internal/linkutil"
)

func rangeFromRequest(c *gin.Context) string {
	if r := c.Query("range"); r != "" {
		return r
	}
	return c.GetHeader("Range")
}

func parseListRange(c *gin.Context) (offset, limit, rangeStart, rangeEnd int, hasRange bool, err error) {
	rangeRaw := rangeFromRequest(c)
	if rangeRaw == "" {
		return 0, 0, 0, 0, false, nil
	}
	rangeStart, rangeEnd, err = linkutil.ParseRange(rangeRaw)
	if err != nil {
		return 0, 0, 0, 0, false, err
	}
	return rangeStart, rangeEnd - rangeStart, rangeStart, rangeEnd, true, nil
}

func writeRangeError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range"})
}
