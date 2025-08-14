package main

import (
	"context"
	"fmt"

	pagination "github.com/gin-melodic/gf-pagination"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	// Create paginator instance
	paginator := pagination.New(
		pagination.WithDefaultPageSize(10),
		pagination.WithMaxPageSize(100),
	)

	// Pagination request
	req := &pagination.Request{
		Page:     1,
		PageSize: 20,
		OrderBy:  "id",
		Sort:     "desc",
	}

	// Query configuration
	query := &pagination.Query{
		Table:  "users",
		Fields: "id,username,email,created_at",
		Conditions: map[string]interface{}{
			"status":        1,
			"username_like": "admin",
		},
	}

	// Execute pagination query
	result, err := paginator.Paginate(context.Background(), req, query)
	if err != nil {
		g.Log().Fatal(context.Background(), err)
	}

	fmt.Printf("Total records: %d\n", result.Total)
	fmt.Printf("Current page: %d\n", result.Page)
	fmt.Printf("Page size: %d\n", result.PageSize)
	fmt.Printf("Total pages: %d\n", result.TotalPages)
	fmt.Printf("Has next page: %v\n", result.HasNext)
}