package management

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/usage"
)

// GetUsageLogs returns paginated, filterable request log entries from SQLite.
func (h *Handler) GetUsageLogs(c *gin.Context) {
	params := usage.LogQueryParams{
		Page:   intQueryDefault(c, "page", 1),
		Size:   intQueryDefault(c, "size", 50),
		Days:   intQueryDefault(c, "days", 7),
		APIKey: strings.TrimSpace(c.Query("api_key")),
		Model:  strings.TrimSpace(c.Query("model")),
		Status: strings.TrimSpace(c.Query("status")),
	}

	result, err := usage.QueryLogs(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filters, err := usage.QueryFilters(params.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats, err := usage.QueryStats(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":   result.Items,
		"total":   result.Total,
		"page":    result.Page,
		"size":    result.Size,
		"filters": filters,
		"stats":   stats,
	})
}

func intQueryDefault(c *gin.Context, key string, def int) int {
	v := strings.TrimSpace(c.Query(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return def
	}
	return n
}
