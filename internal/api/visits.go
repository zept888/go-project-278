package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zept888/go-project-278/internal/linkutil"
	"github.com/zept888/go-project-278/internal/store"
)

func (h *Links) ListVisits(c *gin.Context) {
	ctx := c.Request.Context()

	offset, limit, rangeStart, rangeEnd, hasRange, err := parseListRange(c)
	if err != nil {
		writeRangeError(c)
		return
	}

	total, err := h.Store.CountVisits(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasRange {
		offset, limit = 0, int(total)
	}

	visits, err := h.Store.ListVisits(ctx, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if hasRange {
		c.Header("Content-Range", linkutil.ContentRangeResource("link_visits", rangeStart, rangeEnd, total))
	}

	out := make([]map[string]any, len(visits))
	for i, v := range visits {
		out[i] = visitToResponse(v)
	}
	c.JSON(http.StatusOK, out)
}

func visitToResponse(v store.LinkVisit) map[string]any {
	return map[string]any{
		"id":         v.ID,
		"link_id":    v.LinkID,
		"created_at": v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"ip":         v.IP,
		"user_agent": v.UserAgent,
		"status":     v.Status,
	}
}
