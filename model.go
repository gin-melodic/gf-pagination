package pagination

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// Request Pagination request structure
type Request struct {
	Page     int    `json:"page" v:"min:1#Page number minimum is 1" dc:"Page number, starting from 1"`
	PageSize int    `json:"page_size" v:"min:1|max:500#Page size minimum is 1|Page size maximum is 500" dc:"Page size"`
	OrderBy  string `json:"order_by" dc:"Sort field"`
	Sort     string `json:"sort" v:"in:asc,desc#Sort method can only be asc or desc" dc:"Sort method: asc for ascending, desc for descending"`
}

// Response Pagination response structure
type Response struct {
	List       interface{} `json:"list" dc:"Data list"`
	Page       int         `json:"page" dc:"Current page number"`
	PageSize   int         `json:"page_size" dc:"Page size"`
	Total      int64       `json:"total" dc:"Total records"`
	TotalPages int64       `json:"total_pages" dc:"Total pages"`
	HasNext    bool        `json:"has_next" dc:"Whether there is a next page"`
	HasPrev    bool        `json:"has_prev" dc:"Whether there is a previous page"`
}

// Query Pagination query configuration
type Query struct {
	Model      *gdb.Model             `json:"-" dc:"Database model"`
	Table      string                 `json:"table" dc:"Table name"`
	Fields     string                 `json:"fields" dc:"Query fields"`
	Conditions map[string]interface{} `json:"conditions" dc:"Query conditions"`
	Joins      []string               `json:"joins" dc:"Join queries"`
	GroupBy    string                 `json:"group_by" dc:"Group by field"`
	Having     string                 `json:"having" dc:"Having condition"`
}

// SetDefaults Set default pagination parameters
func (r *Request) SetDefaults(config *Config) *Request {
	if r.Page <= 0 {
		r.Page = config.DefaultPage
	}
	if r.PageSize <= 0 {
		r.PageSize = config.DefaultPageSize
	}
	if r.PageSize > config.MaxPageSize {
		r.PageSize = config.MaxPageSize
	}
	if r.Sort == "" {
		r.Sort = config.DefaultSort
	}
	if r.OrderBy == "" {
		r.OrderBy = config.DefaultOrderBy
	}
	return r
}

// GetOffset Calculate offset
func (r *Request) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}
