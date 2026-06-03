package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zept888/go-project-278/internal/linkutil"
	"github.com/zept888/go-project-278/internal/store"
)

type Links struct {
	Store   store.Store
	BaseURL string
}

type linkBody struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name"`
}

func (h *Links) List(c *gin.Context) {
	ctx := c.Request.Context()
	rangeRaw := c.Query("range")

	var offset, limit int
	var rangeStart, rangeEnd int
	hasRange := rangeRaw != ""
	if hasRange {
		var err error
		rangeStart, rangeEnd, err = linkutil.ParseRange(rangeRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range"})
			return
		}
		offset, limit = rangeStart, rangeEnd-rangeStart
	}

	total, err := h.Store.Count(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasRange {
		offset, limit = 0, int(total)
	}

	links, err := h.Store.List(ctx, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if hasRange {
		c.Header("Content-Range", linkutil.ContentRange(rangeStart, rangeEnd, total))
	}

	out := make([]map[string]any, len(links))
	for i, link := range links {
		out[i] = linkutil.ToResponse(link.ID, link.OriginalURL, link.ShortName, h.BaseURL)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Links) Create(c *gin.Context) {
	var body linkBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.createLink(c.Request.Context(), body.OriginalURL, strings.TrimSpace(body.ShortName))
	if err != nil {
		writeStoreError(c, err)
		return
	}
	c.JSON(http.StatusCreated, linkutil.ToResponse(link.ID, link.OriginalURL, link.ShortName, h.BaseURL))
}

func (h *Links) createLink(ctx context.Context, originalURL, shortName string) (store.Link, error) {
	if shortName != "" {
		return h.Store.Create(ctx, originalURL, shortName)
	}
	for range 8 {
		name, err := linkutil.RandomShortName(8)
		if err != nil {
			return store.Link{}, err
		}
		link, err := h.Store.Create(ctx, originalURL, name)
		if err == nil {
			return link, nil
		}
		if !errors.Is(err, store.ErrConflict) {
			return store.Link{}, err
		}
	}
	return store.Link{}, store.ErrConflict
}

func (h *Links) Get(c *gin.Context) {
	id, err := linkutil.ParseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	link, err := h.Store.Get(c.Request.Context(), id)
	if err != nil {
		writeStoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, linkutil.ToResponse(link.ID, link.OriginalURL, link.ShortName, h.BaseURL))
}

func (h *Links) Update(c *gin.Context) {
	id, err := linkutil.ParseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body linkBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shortName := strings.TrimSpace(body.ShortName)
	if shortName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short_name is required"})
		return
	}

	link, err := h.Store.Update(c.Request.Context(), id, body.OriginalURL, shortName)
	if err != nil {
		writeStoreError(c, err)
		return
	}
	c.JSON(http.StatusOK, linkutil.ToResponse(link.ID, link.OriginalURL, link.ShortName, h.BaseURL))
}

func (h *Links) Delete(c *gin.Context) {
	id, err := linkutil.ParseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Store.Delete(c.Request.Context(), id); err != nil {
		writeStoreError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Links) Redirect(c *gin.Context) {
	link, err := h.Store.GetByShortName(c.Request.Context(), c.Param("shortName"))
	if err != nil {
		writeStoreError(c, err)
		return
	}
	c.Redirect(http.StatusFound, link.OriginalURL)
}

func writeStoreError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
	case errors.Is(err, store.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "short_name already exists"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
