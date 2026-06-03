package api

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "" || name == "-" {
				return fld.Name
			}
			return name
		})
	}
}

type createLinkPayload struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	ShortName   string `json:"short_name" binding:"omitempty,min=3,max=32"`
}

type updateLinkPayload struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	ShortName   string `json:"short_name" binding:"required,min=3,max=32"`
}

func bindCreatePayload(c *gin.Context) (createLinkPayload, bool) {
	var body createLinkPayload
	if !bindJSON(c, &body) {
		return body, false
	}
	body.ShortName = strings.TrimSpace(body.ShortName)
	return body, true
}

func bindUpdatePayload(c *gin.Context) (updateLinkPayload, bool) {
	var body updateLinkPayload
	if !bindJSON(c, &body) {
		return body, false
	}
	body.ShortName = strings.TrimSpace(body.ShortName)
	return body, true
}

func bindJSON(c *gin.Context, dest any) bool {
	if err := c.ShouldBindJSON(dest); err != nil {
		if writeBindError(c, err) {
			return false
		}
	}
	return true
}

func writeBindError(c *gin.Context, err error) bool {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		writeValidationErrors(c, validationErrors)
		return true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	return true
}

func writeValidationErrors(c *gin.Context, errs validator.ValidationErrors) {
	out := make(map[string]string, len(errs))
	for _, e := range errs {
		out[e.Field()] = e.Error()
	}
	c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": out})
}

func writeConflictError(c *gin.Context, field, message string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"errors": gin.H{field: message},
	})
}
