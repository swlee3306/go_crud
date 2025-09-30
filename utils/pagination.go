package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
}

type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, ne, gt, gte, lt, lte, like, in, not_in
	Value    interface{} `json:"value"`
}

func ParsePagination(c *gin.Context) (int, int, string, string) {
	page := 1
	perPage := 10
	sort := "id"
	order := "asc"
	
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}
	
	if s := c.Query("sort"); s != "" {
		sort = s
	}
	
	if o := c.Query("order"); o != "" {
		if o == "desc" || o == "asc" {
			order = o
		}
	}
	
	return page, perPage, sort, order
}

func ParseFilters(c *gin.Context) []Filter {
	var filters []Filter
	
	// Parse filter parameters
	// Format: filter[field][operator]=value
	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			// Extract field and operator from key like "filter[username][eq]"
			parts := strings.Split(key[7:len(key)-1], "][")
			if len(parts) == 2 {
				field := parts[0]
				operator := parts[1]
				
				if len(values) > 0 {
					filters = append(filters, Filter{
						Field:    field,
						Operator: operator,
						Value:    values[0],
					})
				}
			}
		}
	}
	
	return filters
}

func ApplyPagination(db *gorm.DB, page, perPage int) *gorm.DB {
	offset := (page - 1) * perPage
	return db.Offset(offset).Limit(perPage)
}

func ApplySorting(db *gorm.DB, sort, order string) *gorm.DB {
	if sort != "" {
		return db.Order(fmt.Sprintf("%s %s", sort, order))
	}
	return db
}

func ApplyFilters(db *gorm.DB, filters []Filter) *gorm.DB {
	for _, filter := range filters {
		switch filter.Operator {
		case "eq":
			db = db.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
		case "ne":
			db = db.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
		case "gt":
			db = db.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
		case "gte":
			db = db.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
		case "lt":
			db = db.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
		case "lte":
			db = db.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
		case "like":
			db = db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+fmt.Sprintf("%v", filter.Value)+"%")
		case "in":
			if values, ok := filter.Value.([]interface{}); ok {
				db = db.Where(fmt.Sprintf("%s IN ?", filter.Field), values)
			}
		case "not_in":
			if values, ok := filter.Value.([]interface{}); ok {
				db = db.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), values)
			}
		}
	}
	return db
}

func CalculatePagination(page, perPage int, total int64) Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	
	return Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func PaginateQuery(db *gorm.DB, c *gin.Context) (*gorm.DB, Pagination, error) {
	page, perPage, sort, order := ParsePagination(c)
	filters := ParseFilters(c)
	
	// Apply filters
	db = ApplyFilters(db, filters)
	
	// Count total records
	var total int64
	if err := db.Model(db.Statement.Model).Count(&total).Error; err != nil {
		return nil, Pagination{}, err
	}
	
	// Apply sorting
	db = ApplySorting(db, sort, order)
	
	// Apply pagination
	db = ApplyPagination(db, page, perPage)
	
	// Calculate pagination info
	pagination := CalculatePagination(page, perPage, total)
	pagination.Sort = sort
	pagination.Order = order
	
	return db, pagination, nil
}

func CreatePaginatedResponse(data interface{}, pagination Pagination) map[string]interface{} {
	return map[string]interface{}{
		"data":       data,
		"pagination": pagination,
	}
}

// Search functionality
func ApplySearch(db *gorm.DB, searchTerm string, searchFields []string) *gorm.DB {
	if searchTerm == "" || len(searchFields) == 0 {
		return db
	}
	
	var conditions []string
	var args []interface{}
	
	for _, field := range searchFields {
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
		args = append(args, "%"+searchTerm+"%")
	}
	
	if len(conditions) > 0 {
		query := strings.Join(conditions, " OR ")
		db = db.Where(query, args...)
	}
	
	return db
}

// Date range filtering
func ApplyDateRange(db *gorm.DB, field, startDate, endDate string) *gorm.DB {
	if startDate != "" {
		db = db.Where(fmt.Sprintf("%s >= ?", field), startDate)
	}
	if endDate != "" {
		db = db.Where(fmt.Sprintf("%s <= ?", field), endDate)
	}
	return db
}
